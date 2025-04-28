package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/CristianVega28/goserver/server"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type (
	Execution struct {
		Args    []string
		path    string
		port    string
		mode    string
		Restart chan bool
		File    File
		Server  *server.Server
	}
)

func (exec *Execution) WatcherFile() {
	watcher, err := fsnotify.NewWatcher()

	watcher.Add("./api")
	log.Info().Msg("Creando la gorutine")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:

				if !ok {
					return
				}
				// fmt.Println("event: ", event)

				if event.Has(fsnotify.Write) {
					fmt.Println("modified file:", event.Op.String())
					exec.Restart <- true

					// data, _ := exec.File.ExtractData("./api/api.json")
					// exec.Server.GenrateServer(data)

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				fmt.Println("error: ", err)

			case <-exec.Restart:
				exec.Server.Close()

				data, errorExtractData := exec.File.ExtractData("./api/api.json")

				if errorExtractData != nil {
					log.Error().Msg(errorExtractData.Error())
				}

				time.Sleep(1 * time.Second)
				prevServer := server.Server{}
				srv := prevServer.NewServer()

				exec.Server = &srv

				exec.Server.GenrateServer(data)

			}
		}
	}()

	data, errorExtractData := exec.File.ExtractData("./api/api.json")

	if errorExtractData != nil {
		log.Error().Msg(errorExtractData.Error())
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	exec.Server.GenrateServer(data)
}

func (exec *Execution) StaticFile() {
	data, errorExtractData := exec.File.ExtractData(exec.path)

	if errorExtractData != nil {
		log.Error().Msg(errorExtractData.Error())
	}

	exec.Server.GenrateServer(data)
}

func (exec *Execution) ParserArg() {
	rex := regexp.MustCompile(`--\w+=\S+`)
	matches := rex.FindAllString(strings.Join(exec.Args, " "), -1)

	lo.ForEach(matches, func(item string, key int) {
		splitted := strings.Split(item, "=")
		fmt.Println(splitted)
		if len(splitted) == 2 {
			switch splitted[0] {
			case "--port":
				exec.port = splitted[1]
			case "--path":
				exec.path = splitted[1]
			case "--mode":
				exec.mode = splitted[1]
			default:
				exec.mode = "watch"
				exec.path = "./api/api.json"
				exec.port = ":8000"
			}
		}
	})

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	log.Info().Msg("args " + fmt.Sprintf("%v", exec.path))
	log.Info().Msg("args " + fmt.Sprintf("%v", exec.mode))
	log.Info().Msg("args " + fmt.Sprintf("%v", exec.port))
}

func validateArg(arg string, keyword string) bool {
	return true
}

func (exec *Execution) GetMode() bool {
	return exec.mode == "watch"
}

func (exec *Execution) Run() {
	if exec.GetMode() {
		fmt.Println("Modo watcher")
		go exec.WatcherFile()
	} else {
		exec.StaticFile()
	}

}
