package internal

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/scott-mescudi/oxide/pkg/models"
)

type TaurineContext struct {
	Config   *models.Config
	Writer   io.Writer
	mu       sync.Mutex
	msgQueue chan string
	exit     chan struct{}
}

func Init() (ctx *TaurineContext, err error) {
	ctx = &TaurineContext{}
	if _, err := toml.DecodeFile("taurine.toml", &ctx.Config); err != nil {
		return nil, fmt.Errorf("Failed to parse taurine.toml: %v", err)
	}

	ctx.exit = make(chan struct{})
	ctx.msgQueue = make(chan string, 1000)
	ctx.mu = sync.Mutex{}
	ctx.Writer = os.Stdout

	return ctx, nil

}
