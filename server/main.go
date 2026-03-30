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
)

var (
	sigCh = make(chan os.Signal, 1)
)

func main() {
	config.Init()
	log.InitLogger()
	proxy.InitAuditclient()
	go grpc.Init(sigCh)
	go http.Init(sigCh)
	// listen terminage signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh // Block until signal is received
	log.Logger.Infof("Received signal: %v, shutting down...", sig)
}
