// Copyright 2023 Park Zhou <p@ctriple.cn>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	cache  *ssh.Client
	ssherr error
	alive  sync.Once
)

// client returns ssh client and kickoff a keepalive background
// goroutine only once. the client is lazy initialized and cached.
func client() (*ssh.Client, error) {
	if cache == nil {
		cache, ssherr = dial()
	}

	// kickoff keepalive goroutine
	alive.Do(func() {
		go func() {
			tick := time.Tick(3 * time.Second)
			for {
				for ssherr != nil || cache == nil {
					log.Println(ssherr)
					cache, ssherr = dial()
					<-tick
				}
				ssherr = cache.Wait()
			}
		}()
	})

	return cache, ssherr
}

// dial returns a ssh client which can be used to dial from the last
// hop of the tunnel (the last hop must have network connectivity at
// the remote network).
func dial() (client *ssh.Client, err error) {
	log.Println("establishing tunnel connection...")

	for i, hop := range conf.Hops {
		var (
			user    = env(hop.User)
			pass    = env(hop.Pass)
			key     = file(hop.Key)
			address = net.JoinHostPort(hop.Host, hop.Port)

			signer ssh.Signer
		)

		signer, err0 := ssh.ParsePrivateKey([]byte(key))
		if err0 != nil {
			log.Println(err0)
		}

		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		switch i {
		case 0:
			client, err = ssh.Dial("tcp", address, config)
			if err != nil {
				return
			}
		default:
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
	}

	return
}
