package main

import "github.com/mehrdod/purehttpserver/internal/app"

const configPath = "configs"

func main() {
	app.Run(configPath)
}
