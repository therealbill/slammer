package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
			Usage: "Number of concurrenct clients",
		},
		cli.IntFlag{
			Name:  "size,s",
			Value: 1024,
			Usage: "Message size - number of bytes in the message",
		},
		cli.IntFlag{
			Name:  "keys,k",
			Value: 1000,
			Usage: "Number of keys each client will publish",
		},
		cli.StringFlag{
			Name:  "redis-address,r",
			Value: "127.0.0.1:6379",
			Usage: "Redis address to connect to",
		},
		cli.StringFlag{
			Name:  "redis-auth,a",
			Value: "",
			Usage: "Redis password to use",
		},
		cli.StringFlag{
			Name:  "prefix,p",
			Value: "",
			Usage: "Prefix for the key names",
		},
	}
	app.Action = runPublisher
	app.Run(os.Args)
}

func runPublisher(c *cli.Context) {
	msgcount := c.Int("keys")
	msize := c.Int("size")
	runmax := c.Int("clients")
	var buf bytes.Buffer
	for size := 1; size <= msize; size++ {
		buf.WriteString("x")
	}
	message := buf.String()
	for runner := 1; runner <= runmax; runner++ {
		wg.Add(1)
		go setKey(c, message, msgcount, runner, c.String("prefix"))
	}
	wg.Wait()
	log.Printf("Done")

}

func setKey(c *cli.Context, m string, msgcount int, clientid int, prefix string) {
	defer wg.Done()
	config = client.DialConfig{
		Address:  c.String("redis-address"),
		Password: c.String("redis-auth"),
		MaxIdle:  1000,
		Timeout:  30 * time.Second,
	}
	rc, err := client.DialWithConfig(&config)
	if err != nil {
		log.Fatalf("Unable to connect using config: %+v", config)
	}

	err = rc.ClientSetName(fmt.Sprintf("firetruck-%d", clientid))
	if err != nil {
		log.Fatalf("Unable to execute command. Server says '%v'", err)
	}
	for counter := 1; counter <= msgcount; counter++ {
		key := fmt.Sprintf("%s:key-%d-%d", prefix, clientid, counter)
		err := rc.Set(key, m)
		if err != nil {
			log.Printf("[firetruck] Key '%s' not set. Error is '%v'", key, err)
		}
		time.Sleep(time.Millisecond * 5)
	}
}
