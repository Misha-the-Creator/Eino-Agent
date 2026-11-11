package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type AgentSvcInt interface {
	AgentAns(ctx context.Context, mas *adk.Runner, userRequest string) (string, error)
}

type AgentSvc struct {
	// svc AgentSvcInt
	logger *slog.Logger
}

// func NewHndlr(svc AgentSvcInt, logger *slog.Logger) AgentHndlrInt {
// 	return AgentHndlr{svc: svc, logger: logger}
// }

func NewSvc(logger *slog.Logger) AgentSvcInt {
	return &AgentSvc{logger: logger}
}

func (a *AgentSvc) AgentAns(ctx context.Context, mas *adk.Runner, userRequest string) (string, error) {

	iter := mas.Run(ctx, []*schema.Message{
		schema.UserMessage(userRequest),
	}, adk.WithSessionValues(map[string]any{"user": "mikhail"}))

	answer := ""

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		
		if event.Err != nil {
			a.logger.Error("Ошибка при выполнении MAS", "err", event.Err)
			return "", fmt.Errorf("ошибка выполнения MAS: %w", event.Err)
		}

		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		mo := event.Output.MessageOutput
		msg, err := mo.GetMessage() // сам склеит поток, если включен стриминг
		if err != nil {
			continue
		}

		// 1) Модель решила вызвать инструмент / субагента
		for _, tc := range msg.ToolCalls {
			a.logger.Info("TOOL CALL",
				"agent", event.AgentName,
				"path", event.RunPath,
				"tool", tc.Function.Name,
				"args", tc.Function.Arguments)
		}

		// 2) Результат инструмента
		if mo.Role == schema.Tool {
			a.logger.Info("TOOL RESULT",
				"agent", event.AgentName,
				"tool", mo.ToolName,
				"result", msg.Content)
		}

		// 3) Обычный ответ модели
		if mo.Role == schema.Assistant && len(msg.ToolCalls) == 0 {
			a.logger.Info("ANSWER", "agent", event.AgentName, "text", msg.Content)
		}
		}


	return answer, nil
}