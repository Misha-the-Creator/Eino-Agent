package llm

import (
	"context"
	"log/slog"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func BuildCoordinator(ctx context.Context, log *slog.Logger, agents []tool.BaseTool, llmClient *openai.ChatModel) (*adk.ChatModelAgent, error) {

	parentAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "Агент-координатор",
		Description: "Агент-координирующий задачи субагентам",
		Instruction: "Ты — агент-координатор. Твоя задача — перенаправлять пользовательский запрос соответствующему субагенту",
		Model:       llmClient,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agents,
			},
		},
	})

	if err != nil {
		log.Error("Ошибка сборки координатора", "err", err)
		return nil, err
	}

	return parentAgent, nil

}
