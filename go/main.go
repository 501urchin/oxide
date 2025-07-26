package main

import (
	"fmt"
	"log"

	"github.com/501urchin/oxide/pkg/models"
	"github.com/BurntSushi/toml"
)

func main() {
	var conf models.Config

	if _, err := toml.DecodeFile("taurine.toml", &conf); err != nil {
		log.Fatalf("Failed to parse TOML: %v", err)
	}

	fmt.Println(conf)

}
