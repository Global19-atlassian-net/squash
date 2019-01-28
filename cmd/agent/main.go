package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/solo-io/squash/pkg/debuggers"
	"github.com/solo-io/squash/pkg/debuggers/dlv"
	"github.com/solo-io/squash/pkg/debuggers/gdb"
	"github.com/solo-io/squash/pkg/debuggers/java"
	"github.com/solo-io/squash/pkg/debuggers/nodejs"
	"github.com/solo-io/squash/pkg/debuggers/python"

	"github.com/solo-io/squash/pkg/platforms"
	"github.com/solo-io/squash/pkg/platforms/kubernetes"
	"github.com/solo-io/squash/pkg/version"
)

func main() {
	log.SetLevel(log.DebugLevel)

	customFormatter := new(log.TextFormatter)
	log.SetFormatter(customFormatter)

	log.Info(fmt.Sprintf("bridge started %v, %v", version.Version, version.TimeStamp))

	var err error
	var cp platforms.ContainerProcess

	cp, err = kubernetes.NewContainerProcess()
	if err != nil {
		cp, err = kubernetes.NewCRIContainerProcessAlphaV1()
		if err != nil {
			log.WithError(err).Fatal("Cannot get container process locator")
		}
	}

	err = debuggers.RunSquashAgent(getDebugger, cp)
	log.WithError(err).Fatal("Error running debug bridge")

}

func getDebugger(dbgtype string) debuggers.Debugger {
	var g gdb.GdbInterface
	var d dlv.DLV
	var j java.JavaInterface
	var p python.PythonInterface

	switch dbgtype {
	case "dlv":
		return &d
	case "gdb":
		return &g
	case "java":
		return &j
	case "nodejs":
		return nodejs.NewNodeDebugger(nodejs.DebuggerPort)
	case "nodejs8":
		return nodejs.NewNodeDebugger(nodejs.InspectorPort)
	case "python":
		return &p
	default:
		return nil
	}
}
