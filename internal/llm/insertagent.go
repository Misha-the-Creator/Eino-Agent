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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InsertFields struct {
	FileName string `json:"row_info" jsonschema:"required" jsonschema_description:"название файла"`
}

func InsertRow(ctx context.Context, q *generated.Queries, log *slog.Logger) (tool.InvokableTool, error) {
	toolFunc := func(ctx context.Context, input *InsertFields) (string, error) {
		fileNameStr := pgtype.Text{String: input.FileName, Valid: true}
		uuid := pgtype.UUID{Bytes: uuid.New(), Valid: true}
		err := q.InsertRow(ctx, generated.InsertRowParams{Uuid: uuid, Filename: fileNameStr})
		if err != nil {
			log.Error("Произошла ошибка при вставке по названию")
			return "not inserted", err
		}
		return "inserted", err
	}
	return utils.InferTool("insert_by_filename", "Вставить по названию файла", toolFunc)
}

func CreateDBInsertSubagent(ctx context.Context, llm *openai.ChatModel, log *slog.Logger, q *generated.Queries) (tool.BaseTool, error) {

	toolFunc, err := InsertRow(ctx, q, log)

	if err != nil {
		log.Error("Произошла ошибка при вставочке по названию")
		return nil, err
	}

	subAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "InsertAgent",
		Description: "Если нужно внести изменения в БД",
		Instruction: "Ты — ассистент по работе с базой данных и изменяешь ее так, как хочет пользователь",
		Model:       llm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{toolFunc},
		}},
	})

	agentTool := adk.NewAgentTool(ctx, subAgent)

	return agentTool, nil
}
