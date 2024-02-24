package command

import (
	"fmt"
	"github.com/pjimming/mydocker/container"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var ExecCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container, mydocker exec [containerId] [command]",
	Action: func(ctx *cli.Context) error {
		if os.Getenv(container.EnvExecPid) != "" {
			logrus.Infof("pid callback, [pid = %d]", os.Getpid())
			return nil
		}
		// mydocker exec [containerId] [command]
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing containerId or command")
		}
		containerId := ctx.Args().Get(0)
		var cmdArray []string
		cmdArray = append(cmdArray, ctx.Args().Tail()...)
		return execContainer(containerId, cmdArray)
	},
}

func execContainer(containerId string, cmdArray []string) error {
	return container.Exec(containerId, cmdArray)
}
