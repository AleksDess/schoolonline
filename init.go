package main

import (
	"schoolonline/config"
	"schoolonline/tgbot"

	"time"
)

func GoToInit() {}

func init() {
	config.GetConfig()
	<-time.After(300 * time.Millisecond)
	go tgbot.RunBot()
}
