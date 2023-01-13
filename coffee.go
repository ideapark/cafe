// Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// hop represents a hop in the way of an ssh tunnel path, you can use
// either password auth or publickey auth (RFC 4252), the first
// success auth will actually take effect.
type hop struct {
	Host string // ip or dns
	Port string // port number
	User string // ssh user
	Pass string // ssh password auth
	Key  string // ssh publickey auth
}

// coffee describes the global configurations
type coffee struct {
	Wild string   // wild dns suffix resolves to 127.0.0.1
	Urls []string // remote http(s) url to relay
	Hops []hop    // hops of ssh tunnel
}

var (
	port    int
	trace   bool
	version bool

	//go:embed coffee.json
	fs   embed.FS
	conf coffee
	rtls = map[string]bool{}
)

func init() {
	flag.IntVar(&port, "port", 2046, "use another serving port")
	flag.BoolVar(&trace, "trace", true, "trace every http roundtrip object")
	flag.BoolVar(&version, "version", false, "print coffee version")

	log.SetPrefix("🍵 ")

	data, err := fs.ReadFile("coffee.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	for _, raw := range conf.Urls {
		rurl, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("url: %s, error: %s\n", raw, err)
		}
		switch rurl.Scheme {
		case "http":
			rtls[rurl.Host] = false
		case "https":
			rtls[rurl.Host] = true
		default:
			log.Fatalf("%s: scheme [%s] not supported (http or https only)\n", raw, rurl.Scheme)
		}
	}
}

func main() {
	flag.Parse()

	if version {
		fmt.Println(vertag())
		os.Exit(0)
	}

	fmt.Printf(`%s 🍵🍵🍵

The following remote network http(s) are relayed as local http:

%s

`, vertag(), doc())

	doRelay()
}

func doc() string {
	var localurls []string

	for _, raw := range conf.Urls {
		rurl, _ := url.Parse(raw) // err already checked at init stage

		var (
			localscheme  = "http://"
			localhost    = rurl.Host + conf.Wild
			localaddress string
		)

		switch port {
		case 80:
			// default http port, keep it clean
			localaddress = localhost
		default:
			localaddress = net.JoinHostPort(localhost, strconv.Itoa(port))
		}
		localurls = append(localurls, localscheme+localaddress)
	}

	// prepend each local url with emoji
	for i := range localurls {
		localurls[i] = "🛰️ " + localurls[i]
	}

	return strings.Join(localurls, "\n")
}
