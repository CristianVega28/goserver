package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/CristianVega28/goserver/server"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type (
	Execution struct {
		Args          []string
		path          string
		port          string
		mode          string
		File          File
		Server        *server.Server
		MapMiddleware server.MapMiddleware
	}
)

func (exec *Execution) Run() {
	data, errorExtractData := exec.File.ExtractData(exec.path)

	if errorExtractData != nil {
		log.Error().Msg(errorExtractData.Error())
	}

	exec.Server.Srv.Addr = exec.port
	exec.Server.GenrateServer(data)
}

func (exec *Execution) ParserArg() {
	rex := regexp.MustCompile(`--\w+=\S+`)
	fmt.Println(exec.Args)
	fmt.Println(len(exec.Args))
	matches := rex.FindAllString(strings.Join(exec.Args, " "), -1)
	lo.ForEach(matches, func(item string, key int) {
		splitted := strings.Split(item, "=")
		fmt.Println(splitted)
		if len(splitted) == 2 {
			switch splitted[0] {
			case "--port":
				exec.port = fmt.Sprintf(":%s", splitted[1])
			case "--path":
				exec.path = splitted[1]
			case "--mode":
				exec.mode = splitted[1]

			}
		}
	})

	if exec.port == "" {
		exec.port = ":8000"
	}
	if exec.path == "" {
		exec.path = "./api/api.json"
	}
	if exec.mode == "" {
		exec.mode = "static"
	}
}
