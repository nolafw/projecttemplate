package main

import (
	"flag"

	"github.com/nolafw/projecttemplate/internal"
	"github.com/nolafw/projecttemplate/internal/scheduler"
)

// main entry point

func main() {

	env := flag.String("e", "local", "it must be either [local], [develop], [staging], [production] or [testing].")
	isScheduler := flag.Bool("scheduler", false, "if true, run as scheduler.")
	flag.Parse()
	if env == nil {
		panic("Please specify the environment with -e option. It must be either [local], [develop], [staging], [production] or [testing].")
	}

	if isScheduler != nil && *isScheduler {
		scheduler.Start(env)
	} else {
		internal.Run(env)
	}
}
