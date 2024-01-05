package main

import "openai-sdk-experiment/router"

//go:generate swag init -o ./docs
func main() {
	router.Server()
}
