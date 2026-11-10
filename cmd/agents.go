package main

import (
	"context"

	"github.com/Misha-the-Creator/Eino-Agent/generated"
	"github.com/Misha-the-Creator/Eino-Agent/internal/llm"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
)

func BuildMAS(ctx context.Context, appSettings *App, q *generated.Queries) (*adk.Runner, error) {

	log := appSettings.Logger

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

	agentRunner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agentCoordinator, EnableStreaming: true})

	return agentRunner, nil
}
