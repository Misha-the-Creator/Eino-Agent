package agent

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserMsg struct {
	Message string `json:"msg"`
}

type AgentHndlrInt interface {
	AgentRoute(w http.ResponseWriter, r *http.Request)
}

type AgentHndlr struct {
	// svc AgentSvcInt
	logger *slog.Logger
}

// func NewHndlr(svc AgentSvcInt, logger *slog.Logger) AgentHndlrInt {
// 	return AgentHndlr{svc: svc, logger: logger}
// }

func NewHndlr(logger *slog.Logger) AgentHndlrInt {
	return &AgentHndlr{logger: logger}
}

func(a *AgentHndlr) AgentRoute(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var message UserMsg

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		a.logger.Error("ошибка при чтении входа с хэндлера", "err", err)
		return
	}
	a.logger.Info("Что пришло с хэндлера", "msg", message)
}