package agent

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/cloudwego/eino/adk"
)

type UserMsg struct {
	Message string `json:"msg"`
}

type AgentHndlrInt interface {
	AgentRoute(w http.ResponseWriter, r *http.Request)
}

type AgentHndlr struct {
	svc    AgentSvcInt
	logger *slog.Logger
	mas    *adk.Runner
}

// func NewHndlr(svc AgentSvcInt, logger *slog.Logger) AgentHndlrInt {
// 	return AgentHndlr{svc: svc, logger: logger}
// }

func NewHndlr(svc AgentSvcInt, logger *slog.Logger, mas *adk.Runner) AgentHndlrInt {
	return &AgentHndlr{svc: svc, logger: logger, mas: mas}
}

func (a *AgentHndlr) AgentRoute(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var message UserMsg

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		a.logger.Error("ошибка при чтении входа с хэндлера", "err", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	a.logger.Info("Что пришло с хэндлера", "msg", message)

	modelAns, err := a.svc.AgentAns(a.mas, message.Message, r.Context())

	if err != nil {
		a.logger.Error("Ошибка при получении ответа от MAS", "err", err)
		http.Error(w, "MAS unavailable", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(modelAns); err != nil {
		log.Printf("failed to encode response: %v", err)
		return
	}

}
