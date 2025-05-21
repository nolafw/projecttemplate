package util

import (
	"fmt"
	"sync"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/config/pkg/configschema"
	"github.com/nolafw/config/pkg/runtimeconfig"
)

var schemas map[string]*configschema.Definition
var params map[string]*runtimeconfig.Parameters
var onceForConfig sync.Once

func LoadConfig(env *string) {
	onceForConfig.Do(func() {
		paths, err := config.ListModulePathsWithConfig("./internal", "config")
		if err != nil {
			panic(err)
		}
		schemas, params, err = config.Load(*env, paths, true)
		if err != nil {
			panic(err)
		}
	})
}

func GetConfigParam(name string) (*runtimeconfig.Parameters, error) {
	if params == nil {
		panic("config not loaded")
	}
	param, ok := params[name]
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", name)
	}
	return param, nil
}

// TODO: Schemaのgetter
