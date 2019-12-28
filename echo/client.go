package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var port = flag.String("port", "9000", "port")

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", ":"+*port)
	panicOnErr(err)
	defer conn.Close()

	done := make(chan string)

	go handleWrite(conn, done)
	go hanleRead(conn, done)

	fmt.Println(<-done)
	fmt.Println(<-done)
}

func handleWrite(conn net.Conn, done chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		wn, err := conn.Write([]byte(msg))
		panicOnErr(err)
		log.Printf("write %d bytes\n", wn)
	}
	done <- "sent"
}

func hanleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	for {
		nr, err := conn.Read(buf)
		if err == io.EOF {
			log.Println(err)
			break
		}
		panicOnErr(err)
		log.Printf("read %d bytes", nr)
		buf = make([]byte, 1024)
	}
	done <- "read"
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}