package main

import (
	"context"
	"log/slog"

	"github.com/Misha-the-Creator/Eino-Agent/generated"
	"github.com/Misha-the-Creator/Eino-Agent/internal/llm"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
)

func BuildMAS(ctx context.Context, appSettings App, log *slog.Logger, q *generated.Queries) (*adk.ChatModelAgent, error) {

	// Инициализация LLM-клиента
	llmClient, err := llm.GetLLM(ctx, appSettings.apiKey,
		appSettings.modelName,
		log)

	if err != nil {
		log.Error("Ошибка при сборке LLM-клиента", "err", err)
		return nil, err
	}
	//

	// Построение субагентов
	dbSubagent, err := llm.CreateDBSubagent(ctx, llmClient, log, q)

	if err != nil {
		log.Error("Ошибка при сборке MAS", "err", err)
		return nil, err
	}

	subagentsList := []tool.BaseTool{dbSubagent}

	agentCoordinator, err := llm.BuildCoordinator(ctx, log, subagentsList, llmClient)

	if err != nil {
		log.Error("Ошибка при сборке MAS", "err", err)
		return nil, err
	}

	return agentCoordinator, nil
}
