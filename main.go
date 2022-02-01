// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var port = "8080"

func wrong(err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("Error: %v\n", err)
	return true
}

func main() {
	flag.StringVar(&port, "p", port, "Localhost port to listen on")
	flag.Parse()

	// Start the server and listen for incoming connections.
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening for debug requests on localhost:" + port)

	// Close the listener when the application closes.
	defer l.Close()

	// run loop forever, until exit.
	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		// Handle connections concurrently in a new goroutine.
		go handle(c)
	}
}

var one = []byte("1\n")
var zero = []byte("0\n")

func handle(conn net.Conn) {
	defer conn.Close()
	// Buffer client input until a newline.
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	if wrong(err) {
		return
	}

	str := strings.TrimSpace(string(buffer))
	colon := strings.Index(str, ":")
	if colon == -1 {
		fmt.Printf("Client request string lacks colon: '%s'\n", str)
		return
	}
	pid := str[:colon]
	bin := str[colon+1:]

	_, err = strconv.Atoi(string(pid))
	if wrong(err) {
		conn.Write(zero)
		return
	}
	conn.Write(one)
	// expect child process has gone to sleep for 10 seconds

	fmt.Printf("About to gdlv attach %s %s\n", pid, bin)

	gdlv := exec.Command("gdlv", "attach", pid, bin)
	stdOutErr, err := gdlv.CombinedOutput()
	if wrong(err) {
		fmt.Printf("gdlv failed:\n%s\n", stdOutErr)
	} else {
		fmt.Printf("Done with gdlv attach %s %s\n", pid, bin)
	}
}
