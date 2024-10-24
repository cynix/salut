package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brutella/dnssd"
	"github.com/brutella/dnssd/log"
)


func main() {
	version := flag.Bool("version", false, "show version and exit")
	cf := flag.String("config", "/usr/local/etc/salut.yaml", "config file")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	var config Config

	if err := config.Load(*cf); err != nil {
		log.Info.Fatal("failed to load config:", err)
	}

	r, err := dnssd.NewResponder()
	if err != nil {
		log.Info.Fatal("failed to create responder:", err)
	}

	ctx, cancel  := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		cancel()
		log.Info.Print("shutting down...")
	}()

	go func() {
		time.Sleep(1 * time.Second)

		for name, node := range config.Nodes {
			for typ, svc := range node.Services {
				select {
				case <-ctx.Done():
					return
				default:
				}

				s, err := dnssd.NewService(dnssd.Config{
					Name: name,
					Type: typ,
					Host: node.Host,
					Text: svc.Text,
					Port: svc.Port,
					Ifaces: node.Interfaces,
				})
				if err != nil {
					log.Info.Fatalf("failed to create service %q for node %q: %v", typ, name, err)
				}

				if _, err = r.Add(s); err != nil {
					log.Info.Fatalf("failed to register service %q for node %q: %v", typ, name, err)
				}

				log.Info.Printf("registered service %q for node %q", typ, name)
			}
		}
	}()

	if err = r.Respond(ctx); err != nil && err != context.Canceled {
		log.Info.Fatal(err)
	}
}

var (
	Version string = "dev"
)
