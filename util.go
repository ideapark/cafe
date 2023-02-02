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

// FIXME(park): this is a copy from go1.20 strings.CutPrefix(), we
// should use that one when go1.20 is released out.
func cutprefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

// env returns env specified by envkey, the exact string will be
// returned if it's not an envkey.
func env(envkey string) string {
	key, ok := cutprefix(envkey, keyEnv)
	if !ok {
		return envkey
	}
	return os.Getenv(key)
}

// file returns data read from the file specified by filekey. the
// exact string will be returned if it's not a filekey.
func file(filekey string) string {
	file, ok := cutprefix(filekey, keyFile)
	if !ok {
		return filekey
	}

	// expand ~/.ssh/id_rsa like path relative to the current
	// user's home directory (when you're running coffee with
	// `sudo`, the user home should be root's home, which may not
	// always be expected and noted).
	if file, ok = cutprefix(file, "~"); ok {
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

	switch {
	case tls0[host]:
		address = net.JoinHostPort(host, "443")
	default:
		address = net.JoinHostPort(host, "80")
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
