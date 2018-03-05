// +build testrunmain

package main

import (
	"os"
	"testing"
	"github.com/Polarishq/middleware/framework/log"
	"flag"
)

var foo string
func init() {
	flag.StringVar(&foo, "host", "", "Refer to server usage information")
	flag.StringVar(&foo, "port", "", "Refer to server usage information")
	flag.StringVar(&foo, "tls-host", "", "Refer to server usage information")
	flag.StringVar(&foo, "tls-port", "", "Refer to server usage information")
	flag.StringVar(&foo, "tls-certificate", "", "Refer to server usage information")
	flag.StringVar(&foo, "tls-key", "", "Refer to server usage information")
	flag.StringVar(&foo, "static", "", "Refer to server usage information")
	flag.Parse()
}

func TestRunMain(t *testing.T) {
	for i, arg := range os.Args {
		log.WithFields(log.Fields{
			"arg": arg,
			"index": i,
		}).Infof("inspecting argument")
	}
	log.Info("Executing server with code coverage enabled")
	main()
}