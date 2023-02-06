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
// exact string will be returned if it's not a filekey.
func file(filekey string) string {
	file, ok := strings.CutPrefix(filekey, keyFile)
	if !ok {
		return filekey
	}

	// expand ~/.ssh/id_rsa like path relative to the current
	// user's home directory (when you're running coffee with
	// `sudo`, the user home should be root's home, which may not
	// always be expected and noted).
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

// addr returns an address dialable at the remote network.
func addr(httphost string) (address string) {
	host := host(httphost)

	// TODO(park): support non default port
	const (
		portHttps = "443"
		portHttp  = "80"
	)
	switch {
	case tls0[host]:
		address = net.JoinHostPort(host, portHttps)
	default:
		address = net.JoinHostPort(host, portHttp)
	}
	return
}

// host strips the http header host wild dns suffix and port part,
// the remaining should be a meanful address at the remote network.
func host(httphost string) (host string) {
	switch i := strings.Index(httphost, coffee0.Wild); {
	case i > 0:
		host = httphost[0:i]
	default:
		host = httphost
	}
	return
}
