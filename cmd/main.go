package main

import (
	"flag"

	"github.com/nolafw/projecttemplate/internal"
)

// main entry point

func main() {

	env := flag.String("e", "local", "it must be either [local], [develop], [staging], [production] or [testing].")
	flag.Parse()
	if env == nil {
		panic("Please specify the environment with -e option. It must be either [local], [develop], [staging], [production] or [testing].")
	}

	internal.Run(env)
}
