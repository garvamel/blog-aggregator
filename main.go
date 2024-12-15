package main

import (
	"fmt"

	"github.com/garvamel/blog-aggregator/internal/config"
)

func main() {

	configFile := config.Read()
	configFile.SetUser("giapoldo")

	configFile = config.Read()

	fmt.Println(configFile.DBUrl)
	fmt.Println(configFile.CurrentUserName)
}
