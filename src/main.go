package main

import (
	"flag"
	"fmt"
	"os"

	"time"

	"github.com/nats-io/nats"
)

func usage() {
	fmt.Println("Nats-io Message responder.\n\tCMD nats-res <subject to listen for> <response message>")
}

func main() {
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	urls := flag.String("s", nats.DefaultURL, fmt.Sprintf("Nats server URLs seperated by commas (Default: %s)", nats.DefaultURL))

	flag.Usage = usage
	flag.Parse()

	nc, err := nats.Connect(*urls)
	if err != nil {
		fmt.Println("Unable to connecrt to nats server:", err.Error())
		os.Exit(1)
	}
	defer nc.Close()

	args := flag.Args()
	if len(args) < 2 {
		usage()
		return
	}

	sub, res := args[0], []byte(args[1])

	fmt.Println("Now Listening on Nats for the subject:", sub)
	nc.Subscribe(sub, func(m *nats.Msg) {
		if err := nc.Publish(m.Reply, res); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received Subject: %s\nReceived Message: %s\n\nReply Subject: %s\nReply Message: %s\n",
			sub, string(m.Data), m.Reply, string(res))
	})
	select {}
}
