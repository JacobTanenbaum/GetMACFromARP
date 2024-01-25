package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	GetMACAddressFromARP()
}

type res struct {
	hwAddr net.HardwareAddr
	err    error
}

func GetMACAddressFromARP() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("please call with two arguments, an IP address and a number of threads\n")
		return
	}
	ip := args[0]
	num_threads, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("the seconds argument for number of threads is not a number %s: %w\n", num_threads, err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	results := make(chan res)
	for i := 0; i < num_threads; i++ {
		go func(c chan res) {
			hwAddr, _, err := Ping(net.ParseIP(ip))
			c <- res{hwAddr, err}
		}(results)
	}

	returned := 0
	for {
		select {
		case recived := <-results:
			returned++
			fmt.Printf("KEYWORD: number of returned goroutines: %d\n", returned)
			fmt.Printf("\thwAddr:%v \t %s\n", recived.hwAddr, recived.err)

			if returned == num_threads {
				return
			}

		}
	}
}
