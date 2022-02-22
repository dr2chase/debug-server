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
	"sync/atomic"
	"time"
)

func wrong(err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("Error %v\n", err)
	return true
}

// TryDebug checks to see if the DEBUG_SERVER environment variable is not empty,
// and if so, uses that port in a call to DoDebug.
func TryDebug() {
	connPort := os.Getenv("DEBUG_SERVER")
	if connPort == "" {
		return
	}
	DoDebug(connPort)
}

var flag int64 // 0 -> 1 (first, only debugger request) -> 2 (requests will return immediately)

// DoDebug attempts, exactly once, to request that a debugger attach to this
// process, sending a request of the form "pid:executable\n" to localhost:port.
// If no port is provided, the default is "8080".
// Any calls executed will the first call is still active will hang until it
// returns.
func DoDebug(port string) {
	// Don't use defer, this runs in an interesting context and who knows what problems that might cause.
	if !atomic.CompareAndSwapInt64(&flag, 0, 1) {
		for atomic.LoadInt64(&flag) == 1 { // spin, not too busily.
			time.Sleep(100 * time.Millisecond)
		}
		return
	}

	if port == "" {
		port = "8080"
	}

	conn, err := net.Dial("tcp", "localhost:"+port)
	if wrong(err) {
		atomic.StoreInt64(&flag, 2)
		return
	}

	// get the binary name
	bin := os.Args[0]
	if !filepath.IsAbs(bin) {
		if strings.Contains(bin, string(filepath.Separator)) {
			// current directory relative
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Getwd failed %v\n", err)
				conn.Close()
				atomic.StoreInt64(&flag, 2)
				return
			}
			bin = filepath.Join(wd, bin)
		} else {
			bin, err = exec.LookPath(bin)
			if err != nil {
				fmt.Printf("LookPath failed %v\n", err)
				conn.Close()
				atomic.StoreInt64(&flag, 2)
				return
			}
		}
	}

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("Requesting debug of pid:binary %s:%s\n", pid, bin)
	conn.Write([]byte(pid + ":" + bin + "\n"))
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if wrong(err) {
		conn.Close()
		atomic.StoreInt64(&flag, 2)
	}
	if strings.TrimSpace(reply) != "1" {
		conn.Close()
		atomic.StoreInt64(&flag, 2)
		return
	}
	time.Sleep(10 * time.Second) // expect to be interrupted by a debugger here.
	conn.Close()
	atomic.StoreInt64(&flag, 2)
}
