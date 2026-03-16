package main

import (
	"fmt"

	"github.com/kalshiev/gator-go/internal/config"
)

func main() {

	config := config.Read()
	fmt.Println(config)

	config.SetUser("Kalshiev")

	fmt.Println(config)

}
