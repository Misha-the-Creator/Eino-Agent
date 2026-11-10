package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type AgentSvcInt interface {
	AgentAns(mas *adk.Runner, userRequest string, ctx context.Context) (string, error)
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

func (a *AgentSvc) AgentAns(mas *adk.Runner, userRequest string, ctx context.Context) (string, error) {

	iter := mas.Run(ctx, []*schema.Message{
		schema.UserMessage(userRequest),
	}, adk.WithSessionValues(map[string]any{"user": "alice"}))

	answer := ""

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			a.logger.Error("Ошибка при прогоне MAS", "err", event.Err)
			return "", event.Err
		}

		// Событие всегда забираем целиком: в стриминговом режиме Message == nil,
		// и данные лежат в MessageStream.
		var msg *schema.Message
		var err error
		if event.Output != nil && event.Output.MessageOutput != nil && event.Output.MessageOutput.MessageStream != nil {
			st := event.Output.MessageOutput.MessageStream
			var parts []*schema.Message
			for i := 0; ; {
				ch, chErr := st.Recv()
				if chErr != nil {
					err = chErr
					break
				}
				i++
				a.logger.Info("DEBUG chunk", "i", i+1, "chunk", fmt.Sprintf("%+v", ch))
				parts = append(parts, ch)
			}
			st.Close()
			if err != nil && errors.Is(err, io.EOF) {
				err = nil
			}
			if len(parts) > 0 {
				msg, err = schema.ConcatMessages(parts)
			}
		} else {
			msg, err = messageFromEvent(event)
		}
		a.logger.Info("DEBUG событие",
			"agent", event.AgentName,
			"outputNil", event.Output == nil,
			"msgOutNil", event.Output == nil || event.Output.MessageOutput == nil,
			"mvRole", func() any { if event.Output != nil && event.Output.MessageOutput != nil { return event.Output.MessageOutput.Role }; return nil }(),
			"streaming", func() bool { if event.Output != nil && event.Output.MessageOutput != nil { return event.Output.MessageOutput.IsStreaming }; return false }(),
			"err", err,
			"role", func() any { if msg != nil { return msg.Role }; return nil }(),
			"len", func() int { if msg != nil { return len(msg.Content) }; return 0 }(),
			"content", func() string { if msg != nil && len(msg.Content) > 60 { return msg.Content[:60] }; if msg != nil { return msg.Content }; return "" }(),
			"toolCalls", func() int { if msg != nil { return len(msg.ToolCalls) }; return 0 }(),
		)
		if err != nil {
			a.logger.Error("Ошибка при разборе события MAS", "err", err)
			return "", err
		}
		if msg == nil || msg.Role != schema.Assistant || msg.Content == "" {
			continue
		}

		answer = msg.Content
	}

	if answer == "" {
		a.logger.Error("MAS не вернул ответа")
		return "", fmt.Errorf("MAS не вернул ответа")
	}

	return answer, nil
}

// messageFromEvent достаёт сообщение из события раннера,
// склеивая чанки, если ответ пришёл стримом.
func messageFromEvent(event *adk.AgentEvent) (*schema.Message, error) {
	if event == nil || event.Output == nil || event.Output.MessageOutput == nil {
		return nil, nil
	}

	output := event.Output.MessageOutput

	if output.Message != nil {
		return output.Message, nil
	}

	if output.MessageStream == nil {
		return nil, nil
	}

	stream := output.MessageStream
	defer stream.Close()

	return schema.ConcatMessageStream(stream)
}
