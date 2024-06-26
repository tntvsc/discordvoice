package main

import (
	"botpull/configs"
	"botpull/modules/server"
)

func main() {

	cfg := configs.NewConfig("./.env")
	server.NewDiscordServer(cfg).Start()
}
