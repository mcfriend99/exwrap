package main

import (
	"flag"

	"github.com/mcfriend99/exwrap/impl"
)

func main() {
	var cmd impl.CommandLine
	flag.StringVar(&cmd.ConfigFile, "config", impl.DefaultConfigFile, "The exwrap configuration file")
	flag.StringVar(&cmd.BuildDirectory, "dir", impl.DefaultBuildDirectory, "The exwrap build directory")
	flag.Parse()

	// load config file
	config := impl.LoadConfig(cmd)

	impl.Generate(config, cmd)
}
