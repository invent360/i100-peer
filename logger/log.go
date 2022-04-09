package logger

import "github.com/ipfs/go-log/v2"

var logger = log.Logger("rendezvous")

func init() {
	log.SetAllLoggers(log.LevelWarn)
	if err := log.SetLogLevel("rendezvous", "info"); err != nil {
		return
	}
	logger.Info("")
}
