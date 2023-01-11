package main

import (
	"contract_notify/server"
	"contract_notify/config"
)

func main() {
	cfg := config.Defaults
	s, e := server.New(cfg)
	if e != nil {
		panic(e)
	}
	s.Start()
}