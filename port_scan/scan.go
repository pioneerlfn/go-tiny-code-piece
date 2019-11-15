package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func IsOpen(host string, port int, timeout time.Duration) bool {
	// fmt.Printf("scanning %s:%d...\n", host, port)
	time.Sleep(time.Microsecond * 200)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		// fmt.Printf("========= %d OPEN==========\n", port)
		_ = conn.Close()
		return true
	}
	return false
}

func main() {
	var wg sync.WaitGroup
	hostname := flag.String("hostname", "", "host to scan")
	startPort := flag.Int("start-port", 80, "the port on which scan start")
	endPort := flag.Int("end-port", 1000, "the port on which scan end")
	timeout := flag.Duration("timeout", time.Second, "timeout")
	flag.Parse()

	var ports []int
	var mutex sync.Mutex
	for port := *startPort; port < *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if open := IsOpen(*hostname, p, *timeout); open {
				mutex.Lock()
				ports = append(ports, p)
				mutex.Unlock()
			}

		}(port)
	}
	wg.Wait()
	fmt.Printf("========= %s ===========\n", *hostname)
	fmt.Printf("open ports: %v\n", ports)
}
