// Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	keyEnv  = "ENV:"
	keyFile = "FILE:"
)

// env returns env specified by envkey, the exact string will be
// returned if it's not an envkey.
func env(envkey string) string {
	key, ok := strings.CutPrefix(envkey, keyEnv)
	if !ok {
		return envkey
	}
	return os.Getenv(key)
}

// file returns data read from the file specified by filekey. the
// exact string will be returned if it's not an filekey.
func file(filekey string) string {
	file, ok := strings.CutPrefix(filekey, keyFile)
	if !ok {
		return filekey
	}

	// expand ~/.ssh/id_rsa like path relative to the current
	// user's homedir
	if file, ok = strings.CutPrefix(file, "~"); ok {
		home, _ := os.UserHomeDir()
		file = filepath.Join(home, file)
		file = filepath.Clean(file)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("filekey: %s, %s\n", filekey, err)
	}
	return string(data)
}

// raddr returns an address dialable at the remote network.
func raddr(host string) (address string) {
	host = rhost(host)

	switch {
	case rtls[host]:
		address = net.JoinHostPort(host, "443")
	default:
		address = net.JoinHostPort(host, "80")
	}

	return
}

// rhost strips the localaddr wild dns suffix and port part, the
// remaining should be a meanful dns at the remote network.
func rhost(localaddr string) (host string) {
	switch i := strings.Index(localaddr, conf.Wild); {
	case i > 0:
		host = localaddr[0:i]
	default:
		host = localaddr
	}
	return
}
