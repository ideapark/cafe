// Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"text/tabwriter"
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
	view    bool
	conf    string

	//go:embed coffee.json
	fs      embed.FS
	coffee0 coffee
	tls0    = map[string]bool{}
)

func init() {
	flag.IntVar(&port, "port", 2046, "use another serving port")
	flag.BoolVar(&trace, "trace", true, "trace every http roundtrip object")
	flag.BoolVar(&version, "version", false, "print coffee version")
	flag.BoolVar(&view, "view", false, "print default coffee.json")
	flag.StringVar(&conf, "conf", "", "filepath to coffee.json")
}

func main() {
	flag.Parse()

	switch {
	case version:
		fmt.Println(vertag())
		os.Exit(0)
	case view:
		data, _ := fs.ReadFile("coffee.json")
		fmt.Println(string(data))
		os.Exit(0)
	}

	var (
		data []byte
		err  error
	)

	switch {
	case len(conf) > 0:
		data, err = os.ReadFile(conf)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		data, err = fs.ReadFile("coffee.json")
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = json.Unmarshal(data, &coffee0)
	if err != nil {
		log.Fatalln(err)
	}

	for _, raw := range coffee0.Urls {
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("url: %s, error: %s\n", raw, err)
		}
		switch u.Scheme {
		case "http":
			tls0[u.Host] = false
		case "https":
			tls0[u.Host] = true
		default:
			log.Fatalf("%s: scheme [%s] not supported (http or https only)\n", raw, u.Scheme)
		}
	}

	help()
	doRelay()
}

func help() {
	fmt.Println(vertag())

	var (
		buffer = &bytes.Buffer{}
		writer = tabwriter.NewWriter(buffer, 0, 0, 1, ' ', tabwriter.Debug)
	)

	fmt.Fprintln(writer, "#relay", "\t", "Remote http(s)", "\t", "Local http")
	fmt.Fprintln(writer, "------", "\t", "--------------", "\t", "----------")

	for i, raw := range coffee0.Urls {
		u, _ := url.Parse(raw) // err already checked at init stage

		var (
			scheme  = "http://"
			host    = u.Host + coffee0.Wild
			address string
		)
		switch port {
		case 80:
			// default http port, keep it clean
			address = host
		default:
			address = net.JoinHostPort(host, strconv.Itoa(port))
		}
		localurl := scheme + address
		fmt.Fprintln(writer, i+1, "\t", raw, "\t", localurl)
	}
	writer.Flush()

	fmt.Println(buffer.String())
}
