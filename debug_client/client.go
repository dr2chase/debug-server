// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug_client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func wrong(err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("Error %v\n", err)
	return true
}

func TryDebug() {
	connPort := os.Getenv("DEBUG_SERVER")
	if connPort == "" {
		return
	}
	DoDebug(connPort)
}

func DoDebug(port string) {
	if port == "" {
		port = "8080"
	}

	conn, err := net.Dial("tcp", "localhost:"+port)
	if wrong(err) {
		return
	}
	defer conn.Close()

	// get the binary name
	bin := os.Args[0]
	if !filepath.IsAbs(bin) {
		if strings.Contains(bin, string(filepath.Separator)) {
			// current directory relative
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Getwd failed %v\n", err)
				return
			}
			bin = filepath.Join(wd, bin)
		} else {
			bin, err = exec.LookPath(bin)
			if err != nil {
				fmt.Printf("LookPath failed %v\n", err)
				return
			}
		}
	}

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("Requesting debug of pid:binary %s:%s\n", pid, bin)
	conn.Write([]byte(pid + ":" + bin + "\n"))
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if wrong(err) {
		return
	}
	if strings.TrimSpace(reply) != "1" {
		return
	}
	time.Sleep(10 * time.Second) // expect to be interrupted by a debugger here.
}
