// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug_client_test

import (
	"fmt"
	"github.com/dr2chase/debug-server/debug_client"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMe(t *testing.T) {
	// This turns out not to work if go test writes the binary to a temporary file.
	debug_client.TryDebug()
}

func TestManual(t *testing.T) {
	if !testing.Verbose() {
		t.Skip()
	}

	bin := os.Args[0]
	var err error
	if !filepath.IsAbs(bin) {
		if strings.Contains(bin, string(filepath.Separator)) {
			// current directory relative
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Getwd failed %v\n", err)
				t.Fail()
			}
			bin = filepath.Join(wd, bin)
		} else {
			bin, err = exec.LookPath(bin)
			if err != nil {
				fmt.Printf("LookPath failed %v\n", err)
				t.Fail()
			}
		}
	}
	fmt.Printf("Our binary:pid is %s:%d\n", os.Args[0], os.Getpid())
	time.Sleep(60 * time.Second)
}
