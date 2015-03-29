package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var servs map[string]*service
var procs map[string]process

type service struct {
	Enabled bool
	Start   []string
}

type process struct {
	file *os.File
	cmd  *exec.Cmd
}

func init() {
	servs = make(map[string]*service)
	procs = make(map[string]process)
}

func serving(name string) (ok bool) {
	_, ok = servs[name]
	return
}

func running(name string) (ok bool) {
	_, ok = procs[name]
	return
}

func create(path string) (f *os.File, err error) {
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return nil, errors.New("failed to remove file")
		}
	}
	f, err = os.Create(path)
	return
}

func command(args []string) (*os.File, *exec.Cmd, error) {
	f, err := create(path + args[0])
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout, cmd.Stderr = f, f
	return f, cmd, nil
}

func add(args []string) error {
	if !serving(args[0]) && !running(args[0]) {
		servs[args[0]] = &service{Enabled: true, Start: args}
		err := save(path + "services")
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("service already exists")
}

func remove(name string) error {
	if serving(name) {
		if running(name) {
			return errors.New("service is currently running")
		}
		delete(servs, name)
		delete(procs, name)
		err := save(path + "services")
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("service does not exist")
}

func enable(name string) error {
	if serving(name) {
		if servs[name].Enabled {
			return errors.New("service already enabled")
		}
		servs[name].Enabled = true
		err := save(path + "services")
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("service does not exist")
}

func disable(name string) error {
	if serving(name) {
		if !servs[name].Enabled {
			return errors.New("service is already disabled")
		}
		servs[name].Enabled = false
		err := save(path + "services")
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("service does not exist")
	}
}

func start(name string) error {
	if serving(name) {
		if !running(name) {
			if _, ok := procs[name]; !ok {
				f, cmd, err := command(servs[name].Start)
				if err != nil {
					return err
				}
				procs[name] = process{file: f, cmd: cmd}
				err = procs[name].cmd.Start()
				if err != nil {
					return errors.New("failed to start process")
				}
				go procs[name].cmd.Wait()
				return nil
			} else {
				return errors.New("process already exists")
			}
		} else {
			return errors.New("service is already running")
		}
	} else {
		return errors.New("service does not exist")
	}
}

func stop(name string) error {
	if serving(name) {
		if running(name) {
			procs[name].cmd.Process.Kill()
			delete(procs, name)
			return nil
		} else {
			return errors.New("service is not running")
		}
	} else {
		return errors.New("service does not exist")
	}
}

func load(path string) error {
	if len(servs) > 0 || len(procs) > 0 {
		return errors.New("services already loaded")
	}
	b, err := ioutil.ReadFile(path)
	if err == nil {
		var data map[string]*service
		err = json.Unmarshal(b, &data)
		if err == nil {
			servs = data
		}
	}
	return err
}

func reload() error {
	err := unload()
	if err == nil {
		err = load(path + "services")
	}
	return err
}

func unload() error {
	interrupt()
	for key, _ := range procs {
		delete(procs, key)
	}
	for key, _ := range servs {
		delete(servs, key)
	}
	if len(procs) == 0 {
		if len(servs) == 0 {
			return nil
		} else {
			return errors.New("failed to unload all services")
		}
	} else {
		return errors.New("failed to unload all processes")
	}
}

func save(path string) error {
	f, err := create(path)
	if err == nil {
		js, err := json.MarshalIndent(servs, "", "\t")
		if err == nil {
			_, err = f.Write(js)
		}
	}
	return err
}

func interrupt() {
	for _, val := range procs {
		err := val.cmd.Process.Signal(os.Interrupt)
		if err != nil {
			log.Println(err)
		}
	}
}

func boot() error {
	for key, val := range servs {
		if val.Enabled {
			err := start(key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func reboot() error {
	interrupt()
	for key, _ := range procs {
		delete(procs, key)
	}
	err := boot()
	return err
}

func shutdown() {
	interrupt()
	os.Remove(usock)
	os.Exit(0)
}
