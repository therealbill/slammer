package main

import (
	"bytes"
	"log"
	"os"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/therealbill/libredis/client"
)

var (
	config  client.DialConfig
	wg      sync.WaitGroup
	VERSION string
)

func main() {
	app := cli.NewApp()
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "clients,c",
			Value: 50,
			Usage: "Number of concurrent clients",
		},
		cli.IntFlag{
			Name:  "msize,m",
			Value: 1024,
			Usage: "Message size - number of bytes in the message",
		},
		cli.IntFlag{
			Name:  "publish,p",
			Value: 1000,
			Usage: "Number of messages each client will publish",
		},
		cli.StringFlag{
			Name:  "redisaddress,r",
			Value: "localhost:6379",
			Usage: "Redis IP:PORT or DNS:PORT connection string",
		},
		cli.StringFlag{
			Name:  "redisauth,a",
			Usage: "The password to authenticate to the Redis instance",
		},
	}
	app.Action = runPublisher
	app.Run(os.Args)
}

func runPublisher(c *cli.Context) {
	config = client.DialConfig{
		Address:  c.String("redisaddress"),
		Password: c.String("redisauth"),
	}
	rc, err := client.DialWithConfig(&config)
	if err != nil {
		log.Fatalf("Unable to connect using config: %+v", config)
	}

	err = rc.Ping()
	if err != nil {
		log.Fatalf("Unable to execute command. Server says '%v'", err)
	}
	log.Print("Ping successful")
	msgcount := c.Int("publish")
	msize := c.Int("msize")
	runmax := c.Int("clients")
	var buf bytes.Buffer
	for size := 1; size <= msize; size++ {
		buf.WriteString("x")
	}
	message := buf.String()
	for runner := 1; runner <= runmax; runner++ {
		wg.Add(1)
		go publishMessage(rc, message, msgcount)
	}
	wg.Wait()
	log.Printf("Done")

}

func publishMessage(rc *client.Redis, m string, msgcount int) {
	defer wg.Done()
	for counter := 1; counter <= msgcount; counter++ {
		_, err := rc.Publish("test:pubsub", m)
		if err != nil {
			log.Printf("[subscriber] Message not sent. Error is '%v'", err)
		}
	}
}
