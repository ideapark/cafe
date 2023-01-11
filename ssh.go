// Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	cache  *ssh.Client
	ssherr error
	alive  sync.Once
)

// client returns ssh client and kickoff a keepalive background
// goroutine only once. note that the client is lazy initialized and
// cached.
func client() (*ssh.Client, error) {
	if cache == nil {
		cache, ssherr = dial()
	}

	// kickoff keepalive goroutine
	alive.Do(func() {
		go func() {
			for {
				for ssherr != nil {
					log.Println(ssherr)
					cache, ssherr = dial()
				}
				ssherr = cache.Wait()
			}
		}()
	})

	return cache, ssherr
}

// dial returns a ssh client used to dial from the last hop of the
// tunnel (the last hop should have network connectivity at the remote
// network).
func dial() (client *ssh.Client, err error) {
	log.Println("establishing tunnel connection...")

	for _, hop := range conf.Hops {
		hop.User = env(hop.User)
		hop.Pass = env(hop.Pass)
		hop.Key = file(hop.Key)

		var (
			address = net.JoinHostPort(hop.Host, hop.Port)
			signer  ssh.Signer
		)

		signer, err0 := ssh.ParsePrivateKey([]byte(hop.Key))
		if err0 != nil {
			log.Println(err0)
		}

		config := &ssh.ClientConfig{
			User: hop.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(hop.Pass),
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// first hop
		if client == nil {
			client, err = ssh.Dial("tcp", address, config)
			if err != nil {
				return
			}
			continue
		}

		var (
			conn net.Conn
		)
		conn, err = client.Dial("tcp", address)
		if err != nil {
			client.Close()
			return
		}

		var (
			nconn ssh.Conn
			chans <-chan ssh.NewChannel
			reqs  <-chan *ssh.Request
		)
		nconn, chans, reqs, err = ssh.NewClientConn(conn, address, config)
		if err != nil {
			client.Close()
			return
		}

		client = ssh.NewClient(nconn, chans, reqs)
	}

	return
}
