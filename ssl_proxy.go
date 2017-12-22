package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

// SLLProxy proxy actas as an transparent TLS proxy.
func SSLProxy(ctx context.Context) {

	// Load keys
	cert, err := tls.LoadX509KeyPair(flagCertFile, flagKeyFile)
	if err != nil {
		panic(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Start TLS server
	ln, err := tls.Listen("tcp", ":"+flagSSLPort, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = ln.Close()
	}()

	// Accept connections
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			defer func() {
				_ = conn.Close()
			}()
			r := bufio.NewReader(conn)
			for {
				msg, err := r.ReadString('\n')
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Print(msg)
			}
		}(c)
	}
}
