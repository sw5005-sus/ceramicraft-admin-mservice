package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/grpc"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/redis"
)

var (
	sigCh = make(chan os.Signal, 1)
)

func main() {
	config.Init()
	log.InitLogger()
	proxy.InitAuditClient()
	repository.Init()
	log.Logger.Info("Init db completed")
	redis.InitRedis()
	log.Logger.Info("Init redis completed")
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	log.Logger.Infof("Received signal: %v, shutting down...", sig)
}
