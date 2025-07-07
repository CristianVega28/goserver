package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CristianVega28/goserver/core"
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
		MapMiddleware: server.CreateMapMiddleware(),
	}

	exec.ParserArg()
	// var isEvent bool = exec.GetMode()

	// var modelsvar models.Models = models.Models{}
	// modelsvar.Init()
	// models.Migration()
	go func() {
		exec.Run()
	}()
	// Canal para capturar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	fmt.Println("Server PID")
	fmt.Println(os.Getpid())

	go func() {
		<-quit
		log.Println("Apagando servidor...")

		// Cierre con timeout de 5 segundos
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := srv.Srv.Shutdown(ctx); err != nil {
			log.Fatalf("Error cerrando el servidor: %v", err)
		}

		log.Println("Servidor cerrado limpiamente")

	}()
	select {}

}
