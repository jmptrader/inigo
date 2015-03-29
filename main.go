package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
)

const sep = string(os.PathSeparator)

var home, path, usock string

func init() {
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEPATH")
	} else {
		home = os.Getenv("HOME")
	}
	path = home + sep + ".inigo" + sep
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0700)
		check(err)
	}
	usock = path + "server.sock"
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("Got signal:", s)
}

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		data := string(buf[0:n])
		log.Println(string(data))
	}
}

func newTask(tsk string) {
	var t = task{Name: tsk, Args: flag.Args()}
	b, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}
	c := client()
	c.Write(b)
	reader(c)
}

func main() {
	doServ := flag.Bool("server", false, os.Args[0]+" --server")
	doAdd := flag.Bool("add", false, os.Args[0]+" --add <name> <command>")
	doRemove := flag.Bool("remove", false, os.Args[0]+" --remove <name>")
	doEnable := flag.Bool("enable", false, os.Args[0]+" --enable <name>")
	doDisable := flag.Bool("disable", false, os.Args[0]+" --disable <name>")
	doStart := flag.Bool("start", false, os.Args[0]+" --start <name>")
	doStop := flag.Bool("stop", false, os.Args[0]+" --stop <name>")
	doLoad := flag.Bool("load", false, os.Args[0]+" --load <path>")
	doUnload := flag.Bool("unload", false, os.Args[0]+" --unload")
	doReload := flag.Bool("reload", false, os.Args[0]+" --reload")
	doSave := flag.Bool("save", false, os.Args[0]+" --save <path>")
	doReboot := flag.Bool("reboot", false, os.Args[0]+" --reboot")
	doShutdown := flag.Bool("shutdown", false, os.Args[0]+" --shutdown")
	flag.Parse()

	switch {
	case *doServ:
		defer func() {
			shutdown()
		}()
		load(path + "services")
		go hub()
		err := boot()
		check(err)
		wait()
	case *doAdd:
		newTask("add")
	case *doRemove:
		newTask("remove")
	case *doEnable:
		newTask("enable")
	case *doDisable:
		newTask("disable")
	case *doStart:
		newTask("start")
	case *doStop:
		newTask("stop")
	case *doLoad:
		newTask("load")
	case *doUnload:
		newTask("unload")
	case *doReload:
		newTask("reload")
	case *doSave:
		newTask("save")
	case *doReboot:
		newTask("reboot")
	case *doShutdown:
		newTask("shutdown")
	}
}
