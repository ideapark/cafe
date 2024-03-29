// Copyright 2023 Park Zhou <ideapark@139.com>. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// dumpbody returns true if the http header declares that http body is
// human readable (such as json, text, html, css). If the body is
// compressed, it's absolutely not human readable, and false is
// returned.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#directives
func dumpbody(header http.Header) bool {
	switch enc := header.Get("Content-Encoding"); {
	case strings.Contains(enc, "gzip"),
		strings.Contains(enc, "compress"),
		strings.Contains(enc, "deflate"),
		strings.Contains(enc, "br"):
		return false
	}
	// as a bonus, cafe could log all the http roundtrip
	// objects, but for developers only json data will be useful
	// (such as debugging a restful api). other humman readable
	// MIME types will cause too much noise, will just log their
	// header.
	switch ctype := header.Get("Content-Type"); {
	case strings.Contains(ctype, "application/json"):
		return true
	default:
		return false
	}
}

// tuncache caches the http.Transport by address
var (
	tuncache = make(map[string]*http.Transport)
	mu       = &sync.Mutex{}
)

func tunnel(address string) *http.Transport {
	mu.Lock()
	defer mu.Unlock()

	if tun, ok := tuncache[address]; ok {
		return tun
	}

	tun := &http.Transport{
		DialContext: func(_ context.Context, _ string, _ string) (net.Conn, error) {
			client, err := client()
			if err != nil {
				return nil, err
			}
			return client.Dial("tcp", address)
		},
	}
	tuncache[address] = tun
	return tun
}

// every incoming http request will increment it by 1.
var roundno int64

// relay intercepts the http request orginated from the cafe end user:
//  1. establish ssh tunnel to the remote network
//  2. round trip this request
//  3. intercepts the response
//  4. write back the response to the origin requestor
func relay(w http.ResponseWriter, req *http.Request) {
	var (
		roundno = atomic.AddInt64(&roundno, 1)
		scheme  = "http"
		host    = host(req.Host)
		address = addr(req.Host)
	)

	switch {
	case tls0[host]:
		scheme = "https"
	}

	// Setups for relaying this request
	req.URL.Scheme = scheme
	req.URL.Host = address
	req.Host = host

	doTrace(req, roundno)

	// Transport is roundtripping over ssh tunnel connection
	sshtun := tunnel(address)
	defer sshtun.CloseIdleConnections()

	// Setup tls configurations when upstream speaks tls
	if tls0[host] {
		sshtun.TLSClientConfig = &tls.Config{
			ServerName: host,
		}
	}

	resp, err := sshtun.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		doTrace(err, roundno)
		return
	}
	defer resp.Body.Close()

	doTrace(resp, roundno)

	// Copy http headers from the upstream response
	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	// We must write back the status code of the upstream
	// response, or it will always be 200 (StatusOK).
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

// doTrace dumps the http roundtrip objects (aka,. request and
// response) to standard logger. Every roundtrip object pairs will be
// numbered monotonically and incrementally. note that it can also log
// error objects.
func doTrace(obj any, roundno int64) {
	if !trace {
		return
	}

	switch obj := obj.(type) {
	case error:
		log.Printf("#%d\n%v\n", roundno, obj)
	case *http.Request:
		data, err := httputil.DumpRequest(obj, dumpbody(obj.Header))
		if err == nil {
			log.Printf("#%d\n%v\n", roundno, string(data))
		} else {
			log.Printf("#%d\n%v\n", roundno, err)
		}
	case *http.Response:
		data, err := httputil.DumpResponse(obj, dumpbody(obj.Header))
		if err == nil {
			log.Printf("#%d\n%v\n", roundno, string(data))
		} else {
			log.Printf("#%d\n%v\n", roundno, err)
		}
	default:
		log.Fatalln("trace with wrong http object!")
	}
}

func doRelay() {
	http.HandleFunc("/", relay)

	local := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	err := http.ListenAndServe(local, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
