package main

import (
	config "github.com/monsterr00/go-pass-keeper/configs/server"
	"github.com/monsterr00/go-pass-keeper/internal/server/applayer"
	"github.com/monsterr00/go-pass-keeper/internal/server/httplayer"
	"github.com/monsterr00/go-pass-keeper/internal/server/storelayer"
)

func main() {
	config.SetConfig()

	storeLayer := storelayer.New()
	appLayer := applayer.New(storeLayer)
	api := httplayer.New(appLayer)

	api.Engage()
}
