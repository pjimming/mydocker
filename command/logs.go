package command

import (
	"fmt"

	"github.com/pjimming/mydocker/container"

	"github.com/urfave/cli"
)

var LogCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("please input your container id")
		}
		containerId := ctx.Args().Get(0)
		return logContainer(containerId)
	},
}

func logContainer(containerId string) error {
	return container.Log(containerId)
}
