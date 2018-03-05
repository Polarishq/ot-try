package backend

import (
	"github.com/Polarishq/middleware/framework"
	"github.com/Polarishq/middleware/framework/log"
)

//InitializeAll loads environment variables from SSM, configures the DB and loads Avanti config
func InitializeAll() {
	cs, err := framework.NewConfigStore(Constants.ComponentName)
	if err != nil {
		log.Error("Panic while setting up config store: " + err.Error())
		panic(err)
	}
	err = cs.SetupEnvironmentFromSSM()
	if err != nil {
		log.Error("Panic pulling from SSM: " + err.Error())
		panic(err)
	}
}
