package main

import (
	"github.com/rinormaloku/gloo-redirector/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	cmd.ExecuteCmd()
}
