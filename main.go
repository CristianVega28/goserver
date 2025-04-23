package main

import (
	"os"
	"time"

	"github.com/CristianVega28/goserver/core"
	"github.com/CristianVega28/goserver/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	prevServer := server.Server{}
	srv := prevServer.NewServer()

	exec := core.Execution{
		Args:   os.Args,
		File:   core.File{},
		Server: &srv,
	}

	exec.ParserArg()
	// var isEvent bool = exec.GetMode()

	exec.Run()

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// log.Debug().Msg("args " + strings.Join(os.Args, ""))
	log.Info().Msg("Starting server...")

	// if isEvent {
	// 	select {}
	// }
}
