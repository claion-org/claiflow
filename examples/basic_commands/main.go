package main

import (
	"flag"
	"log"

	"github.com/claion-org/claiflow/pkg/client"
)

const APP_NAME = "example-helloworld-client"

func main() {
	server := flag.String("server", "localhost:18099", "flow server url")
	clusterUuid := flag.String("clusteruuid", "", "client's cluster id")
	token := flag.String("token", "", "client's token for server connection")
	loglevel := flag.String("loglevel", "debug", "client's log level. One of: debug(defualt)|info|warn|error")

	flag.Parse()

	log.Printf("%s start\n", APP_NAME)

	if len(*server) == 0 {
		log.Fatalf("flag('server') must exsit\n")
	}

	if len(*clusterUuid) == 0 {
		log.Fatalf("flag('clusteruuid') must exsit\n")
	}

	if len(*token) == 0 {
		log.Fatalf("flag('token') must exsit\n")
	}

	c, err := client.New(client.ClientOptions{
		TargetServer: *server,
		ClusterUuid:  *clusterUuid,
		BearerToken:  *token,
		LogLevel:     *loglevel,
		ConnOptions:  client.ConnOptions{},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s end\n", APP_NAME)
}
