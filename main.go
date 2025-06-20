package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/scott-mescudi/taurine/pkg/models"
)

func main() {
	var conf models.Config
	
	if _, err := toml.DecodeFile("taurine.toml", &conf); err != nil {
		log.Fatalf("Failed to parse TOML: %v", err)
	}

	fmt.Println(conf)

}