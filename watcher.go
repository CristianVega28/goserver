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
		removeCompiledFile()
		reloadByOs()
	}

	PathsI interface {
		setPath(path string)
		getPath() string
	}
	Watcher struct {
		Port string
		Mode string
		Path string
		Cmd  *exec.Cmd
	}

	Paths struct {
		Path        string
		NamePathTmp string
	}
)

var (
	enviroment string        = os.Getenv("GOSERVER_ENV") // Default environment
	pathTmp    string        = "tmp"
	l          utils.LoggerI = &utils.Logger{}
	logs       utils.Logger  = l.Create()
)

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

func (pt *Watcher) removeCompiledFile() {
	pathTmpFull := filepath.Join(pathTmp, "main.exe")
	if _, err := os.Stat(pathTmpFull); !os.IsNotExist(err) {
		err = os.Remove(pathTmpFull)
		if err != nil {
			logs.Fatal(fmt.Sprintf("Error removing compiled file: %s", err.Error()))
		}
	}
}

func (w *Watcher) reloadByOs() {
	if runtime.GOOS == "windows" {
		cmd := w.getCmd()
		pid := cmd.Process.Pid
		pidcurrent := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
		pidcurrent.Stdout = os.Stdout
		pidcurrent.Stderr = os.Stderr
		pidcurrent.Run()
	} else {
		cmd := w.getCmd()
		pid := cmd.Process.Pid
		pidcurrent := exec.Command("kill", fmt.Sprintf("%d", pid))
		pidcurrent.Stdout = os.Stdout
		pidcurrent.Stderr = os.Stderr
		pidcurrent.Run()
	}

}

func main() {
	utils.InitLogger()
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
						w.reloadByOs()
						w.removeCompiledFile()
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
		pathTmpFull := filepath.Join(pathTmp, "main.exe")
		arrayMainFile = []string{"build", "-o", pathTmpFull, "main.go"}
	} else if enviroment == "production" {
		ex, err := os.Executable()
		if err != nil {
			logs.Fatal(fmt.Sprintf("Error getting executable path: %s", err.Error()))
		}
		dir := filepath.Dir(ex)

		binary := filepath.Join(dir, fmt.Sprintf("main-%s-%s", runtime.GOOS, runtime.GOARCH))

		arrayMainFile = []string{binary}
	}

	args := []string{
		fmt.Sprintf("--port=%s", watcher.getPort()),
		fmt.Sprintf("--path=%s", watcher.getPath()),
	}

	if enviroment == "development" {
		build := exec.Command("go", arrayMainFile...)
		err := build.Run()
		if err != nil {
			utils.Log.Fatal(err.Error())
		}
		cmd = exec.Command(filepath.Join(pathTmp, "main.exe"), args...)
	} else if enviroment == "production" {
		cmd = exec.Command(arrayMainFile[0], args...)
	}

	mode := watcher.getMode()
	var errRunner error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Ejecutar y mantenerlo vivo (esperar)
	if mode == "static" {
		errRunner = cmd.Run()
	} else if mode == "watch" {
		errRunner = cmd.Start()
	}

	if errRunner != nil {
		utils.Log.Fatal(errRunner.Error())
	}

	watcher.setCmd(cmd)

	return cmd
}
