package agents

import (
	"context"
	"log/slog"

	"github.com/Misha-the-Creator/Eino-Agent/generated"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
    FileName   string `json:"user_id" jsonschema:"required" jsonschema_description:"Название файла"`
}

func QueryRow(ctx context.Context, q *generated.Queries, log *slog.Logger) (tool.InvokableTool, error) {
	toolFunc := func (ctx context.Context, input *User) (generated.UploadedFile, error) {
		fileNameStr := input.FileName
		row, err := q.GetRow(ctx, pgtype.Text{String: fileNameStr})
		if err != nil {
			log.Error("Произошла ошибка при фетче по названию")
			return generated.UploadedFile{}, err
		}
		return row, nil
	}
	return utils.InferTool("get_row_by_filename", "Получить строку по названию файла", toolFunc)
}

func createDBSubagent(ctx context.Context, key string, modelName string, log *slog.Logger, q *generated.Queries) (*adk.ChatModelAgent, error) {

	llm, err := GetLLM(ctx, key, modelName, log)

	if err != nil {
		log.Error("Ошибка при создании модели выше", "err", err)
	}
	toolFunc, err := QueryRow(ctx, q, log)

	if err != nil {
	log.Error("Произошла ошибка при фетче по названию")
		return nil, err
	}

	subAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name: "Агент по работе с базой данных",
		Description: "Если пользователь требует выполнить поиск наибольших/наименьших значений",
		Instruction: "Ты — ассистент по работе с базой данных",
		Model: llm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{toolFunc},
		}},
	})

	return subAgent, nil
}

func buildCoordinator ()
