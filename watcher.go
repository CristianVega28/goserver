package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/CristianVega28/goserver/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

type (
	WatcherI interface {
		setPort(port string)
		getPort() string
		setMode(mode string)
		getMode() string
		setPath(path string)
		getPath() string
		setCmd(children *exec.Cmd)
		getCmd() *exec.Cmd
		Watch(path string, enviroment string)
	}
	Watcher struct {
		Port string
		Mode string
		Path string
		Cmd  *exec.Cmd
	}
)

var (
	enviroment string        = "development" // Default environment
	pathTmp    string        = ".\\tmp"
	l          utils.LoggerI = &utils.Logger{}
	logs       utils.Logger  = l.Create()
)

func main() {

	_, err := os.Stat(pathTmp)

	if os.IsNotExist(err) {
		os.Mkdir(pathTmp, 0755)
	}

	var watcher WatcherI = &Watcher{}

	mode, path, port := extractMode()

	watcher.setMode(mode)
	watcher.setPath(path)
	watcher.setPort(port)

	// Watcher into as reference
	execution(enviroment, watcher)

	if mode == "watch" {
		logs.Msg(fmt.Sprintf("Watching changes in %s", path))
		watcher.Watch(path, enviroment)
	}

}

func (w *Watcher) Watch(path string, enviroment string) {
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
					if w.getCmd() != nil {
						if runtime.GOOS == "windows" {
							cmd := w.getCmd()
							pid := cmd.Process.Pid
							pidcurrent := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
							pidcurrent.Stdout = os.Stdout
							pidcurrent.Stderr = os.Stderr
							pidcurrent.Run()
						}
					}

					execution(enviroment, w)

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
	<-make(chan struct{})

}

func (w *Watcher) setPort(port string) {
	w.Port = port
}
func (w *Watcher) getPort() string {
	return w.Port
}
func (w *Watcher) setMode(mode string) {
	w.Mode = mode
}
func (w *Watcher) getMode() string {
	return w.Mode
}
func (w *Watcher) setPath(path string) {
	w.Path = path
}
func (w *Watcher) getPath() string {
	return w.Path
}

func (w *Watcher) getCmd() *exec.Cmd {
	return w.Cmd
}

func (w *Watcher) setCmd(children *exec.Cmd) {
	w.Cmd = children
}

func extractMode() (string, string, string) {
	rex := regexp.MustCompile(`--\w+=\S+`)
	var mode string
	var path string
	var port string
	matches := rex.FindAllString(strings.Join(os.Args, " "), -1)
	lo.ForEach(matches, func(item string, key int) {
		splitted := strings.Split(item, "=")
		if len(splitted) == 2 {
			switch splitted[0] {
			case "--mode":
				mode = splitted[1]
			case "--path":
				path = splitted[1]
			case "--port":
				port = splitted[1]
			}
		}
	})

	if port == "" {
		port = "8000"
	}
	if path == "" {
		path = "./api/api.json"
	}
	if mode == "" {
		mode = "watch"
	}

	return mode, path, port
}

func execution(enviroment string, watcher WatcherI) (children *exec.Cmd) {
	var cmd *exec.Cmd
	arrayMainFile := []string{}

	if enviroment == "development" {
		pathTmpFull := filepath.Join(".\\", pathTmp, "main.exe")
		arrayMainFile = []string{"build", "-o", fmt.Sprintf("./%s", pathTmpFull), "main.go"}
	} else if enviroment == "production" {
		arrayMainFile = []string{"./main.exe"} // check out about the os of user
	}

	args := []string{
		fmt.Sprintf("--port=%s", watcher.getPort()),
		fmt.Sprintf("--path=%s", watcher.getPath()),
	}

	fmt.Println(watcher.getMode())
	if enviroment == "development" {
		build := exec.Command("go", arrayMainFile...)
		err := build.Run()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(build.Args)
		cmd = exec.Command(fmt.Sprintf("./%s/main.exe", pathTmp), args...)
	} else if enviroment == "production" {
		cmd = exec.Command(arrayMainFile[0], args...)
	}

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	// }

	// Si quieres enviar entrada también
	// cmd.Stdin = os.Stdin

	mode := watcher.getMode()
	var errRunner error

	// Ejecutar y mantenerlo vivo (esperar)
	if mode == "static" {
		errRunner = cmd.Run()
	} else if mode == "watch" {
		errRunner = cmd.Start()
	}

	if errRunner != nil {
		fmt.Println("Error ejecutando core/main.go:", errRunner)
	}

	watcher.setCmd(cmd)

	return cmd
}
