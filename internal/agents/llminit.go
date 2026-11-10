package agents

import (
	"context"
	"log/slog"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

func GetLLM(ctx context.Context, key string, modelName string, log *slog.Logger) (*openai.ChatModel, error) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{APIKey: key,
	BaseURL: "https://polza.ai/api/v1",
	Model: modelName})

	if err != nil {
		log.Error("произошла ошибка при попытке вызвать LLM через polza.ai", "err", err)
		return nil, err
	}

	return chatModel, nil

}