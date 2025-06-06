package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/CristianVega28/goserver/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

type (
	WatcherI interface {
		Run() error
		Watch(path string) error
	}
	Watcher struct {
	}
)

var l utils.LoggerI = &utils.Logger{}
var logs utils.Logger = l.Create()

func main() {

	fmt.Println("Watcher initialized")
	var watcher WatcherI = Watcher{}
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Si quieres enviar entrada también
	cmd.Stdin = os.Stdin

	mode, path := extractMode()
	// Ejecutar y mantenerlo vivo (esperar)
	err := cmd.Run()

	if mode == "watch" {
		watcher.Watch(path)
	} else {
		watcher.Run()
	}

	if err != nil {
		fmt.Println("Error ejecutando core/main.go:", err)
	}
	cmd.Run()

	<-make(chan struct{})

}

func (w Watcher) Run() error {}
func (w Watcher) Watch(path string) error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		logs.Fatal(err.Error())
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					logs.Msg(fmt.Sprintf("modified -> event: %s, file: %s", event.Op, event.Name))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logs.Fatal(fmt.Sprintf("error: %s", err.Error()))
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		logs.Fatal(fmt.Sprintf("error: %s", err.Error()))
	}

	return nil
}

func extractMode() (string, string) {
	rex := regexp.MustCompile(`--\w+=\S+`)
	var mode string
	var path string
	matches := rex.FindAllString(strings.Join(os.Args, " "), -1)
	lo.ForEach(matches, func(item string, key int) {
		splitted := strings.Split(item, "=")
		if len(splitted) == 2 {
			switch splitted[0] {
			case "--mode":
				mode = splitted[1]
			case "--path":
				path = splitted[1]
			}
		}
	})

	return mode, path
}
