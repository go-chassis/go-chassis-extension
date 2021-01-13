package main

import (
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/examples/schemas"

	_ "github.com/go-chassis/go-chassis-extension/registry/kubernetes"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/kube/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.Hello{}, server.WithSchemaID("HelloService"))
	chassis.RegisterSchema("rest-legacy", &schemas.Legacy{}, server.WithSchemaID("LegacyService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
