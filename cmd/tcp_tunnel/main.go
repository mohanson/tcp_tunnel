package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func link(a, b io.ReadWriteCloser) {
	go func() {
		io.Copy(b, a)
		a.Close()
		b.Close()
	}()
	io.Copy(a, b)
	b.Close()
	a.Close()
}

func move(locale, remote string) {
	ln, err := net.Listen("tcp", locale)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("[%s] Create new forwarding from %s\n", remote, locale)
	for {
		connl, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		connr, err := net.Dial("tcp", remote)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("[%s] Income new connection from %s\n", remote, connl.RemoteAddr())
		go link(connl, connr)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s LOCALE/REMOTE ...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		return
	}
	for _, e := range flag.Args() {
		seps := strings.SplitN(e, "/", 2)
		locale := seps[0]
		remote := seps[1]
		go move(locale, remote)
	}
	select {}
}
