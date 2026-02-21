package models

import (
	"fmt"
	"log/slog"

	"github.com/sid-sun/openwebui-bot/cmd/config"
	"github.com/sid-sun/openwebui-bot/pkg/bot/service"
	tele "gopkg.in/telebot.v3"
)

var logger = slog.Default().With(slog.String("package", "Models"))

func getInlineKeyboardMarkup(currentModel string) (string, [][]tele.InlineButton) {
	var modelOptions [][]tele.InlineButton
	modelInfoMessage := "Here are the available models: \n"

	// Try to fetch models from OpenWebUI API
	models, err := service.FetchModels()
	if err != nil {
		logger.Error("failed to fetch models from API, falling back to config", slog.Any("error", err))
		// Fallback to config-based models
		return getConfigBasedKeyboardMarkup(currentModel)
	}

	// Build keyboard from API response
	for _, model := range models {
		text := model.Name
		if currentModel == model.ID {
			text = fmt.Sprintf("*%s*", model.Name)
		}
		modelOptions = append(modelOptions, []tele.InlineButton{
			{
				Data: "model_" + model.ID,
				Text: text,
			},
		})
	}

	if currentModel == "" {
		currentModel = "default"
	}

	return modelInfoMessage, modelOptions
}

func getConfigBasedKeyboardMarkup(currentModel string) (string, [][]tele.InlineButton) {
	var modelOptions [][]tele.InlineButton
	modelInfoMessage := "Here are the available models (from config): \n"
	if currentModel == "" {
		currentModel = "default"
	}
	for _, modelName := range config.GlobalConfig.ModelNames {
		options := config.GlobalConfig.Models[modelName]
		text := fmt.Sprintf("%s (%s) - %d", modelName, options.Model, options.Tweaks.ContextLength)
		if currentModel == modelName {
			text = fmt.Sprintf("*%s* (%s) - %d", modelName, options.Model, options.Tweaks.ContextLength)
		}
		modelOptions = append(modelOptions, []tele.InlineButton{
			{
				Data: "model_" + modelName,
				Text: text,
			},
		})
	}
	return modelInfoMessage, modelOptions
}
