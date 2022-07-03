package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var options internal.Options

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := gin.Default()

	svc, err := internal.New(ctx, log, options, r)
	if err != nil {
		log.WithError(err).Fatal("internal.New()")
	}

	go func() {
		// service connections
		if err := svc.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	shutdownOnSignal(svc)

	svc.Close()
}

func init() {
	options = internal.Options{
		AppName:           os.Getenv("APP_NAME"),
		ListenAddressHTTP: os.Getenv("LISTEN_ADDRESS_HTTP"),
	}

	log = logrus.WithField("app_name", options.AppName)
}

func waitForShutdownSignal() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until signaled
	sig := <-c

	return sig.String()
}

func shutdownOnSignal(svc *internal.Service) {
	signalName := waitForShutdownSignal()
	log.WithField("signal", signalName).Info("Received signal, starting shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if svc.Shutdown(ctx) {
		log.Info("Shutdown complete")
	} else {
		log.Info("Shutdown timed out")
	}
}
