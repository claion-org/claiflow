package main

import (
	"flag"
	"fmt"
	"log"
	"math"

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

	if err := c.RegisterCommand("helloworld", HelloWorld); err != nil {
		log.Fatal(err)
	}

	if err := c.RegisterCommand("math_pow", MathPow); err != nil {
		log.Fatal(err)
	}

	if err := c.RegisterCommand("swap_command", SwapCommand); err != nil {
		log.Fatal(err)
	}

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s end\n", APP_NAME)
}

// helloworld command
type HelloworldReq struct {
	Name string `json:"name"`
}

func HelloWorld(inputs *HelloworldReq) (string, error) {
	if inputs == nil {
		return "", fmt.Errorf("inputs is nil")
	}
	return fmt.Sprintf("hello %s", inputs.Name), nil
}

// math pow command
type PowReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func MathPow(req *PowReq) (float64, error) {
	if req == nil {
		return 0, fmt.Errorf("req is nil")
	}
	return math.Pow(req.X, req.Y), nil
}

// swap command
type SwapCommandReq struct {
	Param1 interface{} `json:"param1"`
	Param2 interface{} `json:"param2"`
}

type SwapCommandResp struct {
	Value1 interface{} `json:"value1"`
	Value2 interface{} `json:"value2"`
}

func SwapCommand(req *SwapCommandReq) (*SwapCommandResp, error) {
	return &SwapCommandResp{
		Value1: req.Param2,
		Value2: req.Param1,
	}, nil
}
