package llm

import (
	"context"
	"log/slog"

	"github.com/Misha-the-Creator/Eino-Agent/generated"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/jackc/pgx/v5/pgtype"
)

type FileName struct {
	FileName string `json:"file_name" jsonschema:"required" jsonschema_description:"Название файла"`
}

func QueryRow(ctx context.Context, q *generated.Queries, log *slog.Logger) (tool.InvokableTool, error) {
	toolFunc := func(ctx context.Context, input *FileName) (generated.GetRowRow, error) {
		fileNameStr := input.FileName
		row, err := q.GetRow(ctx, pgtype.Text{String: fileNameStr, Valid: true})
		if err != nil {
			log.Error("Произошла ошибка при фетче по названию")
			return generated.GetRowRow{}, err
		}
		return row, nil
	}
	return utils.InferTool("get_row_by_filename", "Получить строку по названию файла", toolFunc)
}

func CreateDBSubagent(ctx context.Context, llm *openai.ChatModel, log *slog.Logger, q *generated.Queries) (tool.BaseTool, error) {

	toolFunc, err := QueryRow(ctx, q, log)

	if err != nil {
		log.Error("Произошла ошибка при фетче по названию")
		return nil, err
	}

	subAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "DataBaseAgent",
		Description: "Если пользователь просит вывести информацию о конкретной строке базы данных",
		Instruction: "Ты — ассистент по работе с базой данных",
		Model:       llm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{toolFunc},
		}},
	})

	agentTool := adk.NewAgentTool(ctx, subAgent)

	return agentTool, nil
}
