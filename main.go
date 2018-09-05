package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gwwfps/assembly-lines/db"
	"github.com/gwwfps/assembly-lines/server"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisUrl      string
	RedisPassword string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	config := Config{}
	err := envconfig.Process("assembly", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	d := db.NewDB(config.RedisUrl, config.RedisPassword)
	s := server.NewServer(d)
	log.Fatal(s.Start())
}
