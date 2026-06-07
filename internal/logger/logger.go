package logger

import (
	"log"
	"os"
)

// Logger is the application's shared logger.
var Logger = log.New(
	os.Stdout,
	"[GARNET] ",
	log.LstdFlags|log.Lshortfile,
)
