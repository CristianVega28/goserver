package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CristianVega28/goserver/core"
	"github.com/CristianVega28/goserver/core/middleware"
	"github.com/CristianVega28/goserver/server"
)

//nolint:unused
func main() {

	prevServer := server.Server{}
	srv := prevServer.NewServer()

	exec := core.Execution{
		Args:          os.Args,
		File:          core.File{},
		Server:        &srv,
		MapMiddleware: middleware.CreateMapMiddleware(),
	}

	exec.ParserArg()
	// var isEvent bool = exec.GetMode()

	// var modelsvar models.Models = models.Models{}
	// modelsvar.Init()
	// models.Migration()
	sign := make(chan os.Signal, 1)

	signal.Notify(sign, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		exec.Run()
	}()
	// Canal para capturar señal de interrupción

	select {
	case <-sign:
		os.Exit(0)
	}
}
