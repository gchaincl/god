package main

import (
//	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

const SOCKET = "god.sock"

func connect(sock string) (net.Conn, error){
	return net.Dial("unix", sock)
}

func listen(sock string) (l net.Listener) {
	l, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}
	return l
}

func shouldStartServer(err error) bool {
	if err == nil {
		return false
	}
	return err.(*net.OpError).Err == syscall.ENOENT
}

func startServerAt(sock string) int {
	log.Println("STARTING SERVER...")
	l := listen(sock)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleRequest(c)
	}
}

func handleRequest(c net.Conn) {
	buf := make([]byte, 512)
	nr, err := c.Read(buf)
	if err != nil {
		log.Println(err)
		return
    }

	resp := executeRequest(string(buf[:nr]))
	if _, err := c.Write([]byte(resp)); err != nil {
		log.Println(err)
		return
	}
}

func executeRequest(cmd string) string {
	log.Printf("Processing: '%s'\n", cmd)
	switch cmd {
	case "LIST":
		return "List"
	default:
		return "unknown"
	}
}

func main() {
	c, err := connect(SOCKET)
	if shouldStartServer(err) {
		os.Exit(startServerAt(SOCKET))
	} else if err != nil {
		panic(err)
	}

	log.Println("CLIENT MODE")
	c.Write([]byte("LIST"))
	buf := make([]byte, 512)
	_, err = c.Read(buf)
	log.Println(string(buf))

}
