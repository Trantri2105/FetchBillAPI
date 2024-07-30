package main

import (
	"backend/initilizer"
	"backend/transport"
)

func main() {
	initilizer.LoadEnvVariable()
	router := transport.NewRouter()
	router.Run()
}
