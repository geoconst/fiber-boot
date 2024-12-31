package main

import "fiber-boot/internal/app"

func main() {
	app.InitLogger()
	server := InitializeServer()
	server.Start()
}
