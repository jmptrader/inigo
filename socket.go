package main

import (
	"encoding/json"
	"log"
	"net"
)

type task struct {
	Name string
	Args []string
}

func handler(c net.Conn, data []byte) {
	log.Println("Server:", string(data))
	var t task
	err := json.Unmarshal(data, &t)
	check(err)
	switch t.Name {
	case "add":
		err := add(t.Args)
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Added " + t.Args[0] + " service"))
	case "remove":
		err := remove(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Removed " + t.Args[0] + " service"))
	case "start":
		err := start(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Started " + t.Args[0] + " service"))
	case "stop":
		err := stop(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Stopped " + t.Args[0] + " service"))
	case "enable":
		err := enable(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Enabled " + t.Args[0] + " service"))
	case "disable":
		err := disable(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Disabled " + t.Args[0] + " service"))
	case "load":
		err := load(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Loaded services from " + t.Args[0]))
	case "unload":
		err := unload()
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Unloaded all services"))
	case "reload":
		err := reload()
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Reloaded all services"))
	case "save":
		err := save(t.Args[0])
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Saved services to " + t.Args[0]))
	case "reboot":
		err := reboot()
		if err != nil {
			log.Println(err)
			c.Write([]byte(err.Error()))
			break
		}
		c.Write([]byte("Rebooted all enabled services"))
	case "shutdown":
		c.Write([]byte("Shutting down server"))
		shutdown()
	}
	c.Close()
}

func listener(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		handler(c, data)
		if err != nil {
			log.Fatal("Write: ", err)
		}
	}
}

func hub() {
	l, err := net.Listen("unix", usock)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go listener(fd)
		}
	}()
}

func client() net.Conn {
	c, err := net.Dial("unix", usock)
	if err != nil {
		panic(err)
	}
	return c
}
