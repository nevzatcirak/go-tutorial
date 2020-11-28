package main

import (
	Config "github.com/nevzatcirak/go-examples/config-server/config"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func init() {
	_ = gotenv.Load()
	Config.InitializeLogger("config-server.log")
	Config.LoadProperties()
}

func main() {
	log.Info(Config.GetAllProperties())
	log.Info(Config.GetProperty("spring.application.name"))
}
