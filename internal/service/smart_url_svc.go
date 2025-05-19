package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	handler "github.com/atyagi9006/smart-url/internal/handlers"
	"go.uber.org/zap"
)

const ServiceName string = "smart-url-svc"

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

func (svc *SmartUrlSVC) StartApiService(ctx context.Context, wg *sync.WaitGroup) {

	svc.setupRoutes()
	wg.Add(1) // avoid race
	go func() {
		defer wg.Done()
		svc.Log.Info("Starting service", zap.String("service", ServiceName))

		go func() {
			http.ListenAndServe(":8080", nil)
		}()

		// Graceful shutdown logic
		<-ctx.Done()
		svc.Log.Info("Stopping API service", zap.String("service", ServiceName))
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

	}()

}

func (svc *SmartUrlSVC) setupRoutes() {
	http.HandleFunc("/shorten", svc.handler.ShortenURLHandler)
	http.HandleFunc("/metrics", svc.handler.MetricsHandler)
	http.HandleFunc("/", svc.handler.RedirectHandler)

}
