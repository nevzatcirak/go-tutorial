package main

import (
	Config "github.com/nevzatcirak/go-examples/config-server/config"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	Config.InitializeLogger("config-server.log")
}

func main() {
	properties := Config.GetProperties()
	log.Info(properties)
}
