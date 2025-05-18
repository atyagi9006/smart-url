package internal

import (
	handler "github.com/atyagi9006/smart-url/internal/handlers"
	"go.uber.org/zap"
)

type SmartUrlSVC struct {
	handler *handler.Handler
	Log     *zap.Logger
}

func NewSmartUrlSVC(log *zap.Logger) *SmartUrlSVC {
	svc := SmartUrlSVC{
		handler: handler.NewHandler(log),
		Log:     log,
	}
	return &svc
}


