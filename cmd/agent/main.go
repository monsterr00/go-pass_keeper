package main

import (
	"flag"

	config "github.com/monsterr00/go-pass-keeper/configs/agent"
	"github.com/monsterr00/go-pass-keeper/internal/agent/applayer"
	"github.com/monsterr00/go-pass-keeper/internal/agent/httplayer"
	"github.com/monsterr00/go-pass-keeper/internal/agent/storelayer"
)

func main() {
	flag.Parse()
	config.SetConfig()

	storeLayer := storelayer.New()
	appLayer := applayer.New(storeLayer)
	api := httplayer.New(appLayer)

	api.Engage()
}
