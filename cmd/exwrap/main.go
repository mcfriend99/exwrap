package main

import (
	"flag"

	"github.com/mcfriend99/exwrap/impl"
)

var OsMatrix []string = []string{"windows", "linux", "darwin"}
var ArchMatrix []string = []string{"windows", "linux", "darwin"}

func main() {
	var cmd impl.CommandLine
	flag.StringVar(&cmd.ConfigFile, "config", impl.DefaultConfigFile, "The exwrap configuration file.")
	flag.StringVar(&cmd.BuildDirectory, "dir", impl.DefaultBuildDirectory, "The exwrap build directory.")
	flag.Parse()

	// load config file
	_ = impl.Generate(impl.LoadConfig(cmd), cmd)
}
