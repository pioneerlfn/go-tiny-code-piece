package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

var port = flag.String("port", "9000", "port")

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", ":"+*port)
	// 先处理err
	panicOnErr(err)
	// 别忘了关闭监听socket
	defer ln.Close()
	log.Printf("listening on %s:port", getAddr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept err: %v", err)
		}

		go echo(conn)
	}

}

func echo(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		nr, err := conn.Read(buf)
		if err == io.EOF {
			log.Println(err)
			return
		}
		panicOnErr(err)
		log.Printf("read %d bytes", nr)

		nw, err := conn.Write(buf[:nr])
		panicOnErr(err)
		log.Printf("write %d bytes", nw)
		buf = make([]byte, 1024)
	}
}



func getAddr() string {
	addrs, err := net.InterfaceAddrs()
	panicOnErr(err)
	var localAddr string
	for _, addr := range addrs {
		if strings.HasPrefix(addr.String(), "10.") {
			localAddr = addr.String()
		}
	}
	return localAddr
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}