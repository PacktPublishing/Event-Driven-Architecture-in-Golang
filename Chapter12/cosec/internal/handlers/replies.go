package handlers

import (
	"eda-in-golang/cosec/internal"
	"eda-in-golang/cosec/internal/models"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/sec"
)

func NewReplyHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*models.CreateOrderData], mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(internal.CreateOrderReplyChannel, handlers, am.GroupName("cosec-replies"))
	return err
}
