package main

import (
	"context"
	"fmt"
	"github.com/StinkyPeach/bridge/adapter/outbound"
	"github.com/StinkyPeach/bridge/log"
	"time"
)

func main() {
	log.Info("this is ssr test")

	ssrNode := make(map[string]interface{})
	ssrNode["name"] = "test"
	ssrNode["type"] = "ssr"
	ssrNode["server"] = "127.0.0.1"
	ssrNode["port"] = 8388
	ssrNode["password"] = "password"
	ssrNode["cipher"] = "rc4-md5"
	ssrNode["obfs"] = "http_post"
	ssrNode["obfs-param"] = ""
	ssrNode["protocol"] = "origin"
	ssrNode["protocol-param"] = ""
	ssrNode["udp"] = true

	p, err := outbound.ParseProxy(ssrNode)

	if err != nil {
		fmt.Println(err)
	}

	ur := "https://www.google.com"
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	latency, err := p.URLTest(ctx, ur)
	if err != nil {
		panic(err)
	}

	log.Info("latency: %d ms", latency)

}
