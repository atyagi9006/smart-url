package start

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/atyagi9006/smart-url/internal/service"
	"go.uber.org/zap"
)

func Run() {

	svc := initSmartSvc()
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	// Wait for interrupt signal to gracefully shutdown the server with wait group
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the api service
	svc.StartApiService(ctx, wg)

	<-quit
	svc.Log.Info("Shutting down " + service.ServiceName)
	// preProcessor.timerRoutineDone <- struct{}{}
	cancelFunc() // Signal cancellation to context.Context
	// Wait for all thread to exit
	wg.Wait()
	// We are done
	svc.Log.Info("Graceful shutdown done for " + service.ServiceName)

}

func initSmartSvc() *service.SmartUrlSVC {
	// Initialize the logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize AwsCloudIntegration
	svc := service.NewSmartUrlSVC(logger)
	return svc
}
