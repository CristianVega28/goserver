package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CristianVega28/goserver/core"
	"github.com/CristianVega28/goserver/server"
	"github.com/CristianVega28/goserver/utils"
)

//nolint:unused
func main() {

	lbk := utils.Logger{}
	log := lbk.Create()
	prevServer := server.Server{}
	srv := prevServer.NewServer()

	exec := core.Execution{
		Args:   os.Args,
		File:   core.File{},
		Server: &srv,
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

	log.Msg("Server running on ->  http://localhost" + exec.GetPort())
	select {
	case <-sign:
		os.Exit(0)
	}
}
