package fazzdb

import (
	"fmt"
	"log"
)

var verbose = MODE_SILENT

// Verbose is a function to enable verbose mode
func Verbose() {
	verbose = MODE_VERBOSE
}

func isVerbose() bool {
	return verbose == MODE_VERBOSE
}

func show(str string) {
	log.Printf("[FazzDB] %s", str)
}

// info NOT FOR USE IN PRODUCTION
func info(str string) {
	if !isVerbose() {
		return
	}

	log.Printf("[FazzDB.Info] %s\n", str)
}

// alert NOT FOR USE IN PRODUCTION
func alert(str string) {
	if !isVerbose() {
		return
	}

	fmt.Printf("[FazzDB.Alert] %s\n", str)
}
