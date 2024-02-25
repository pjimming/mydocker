package command

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/pjimming/mydocker/container"
)

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "Commit container to image",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("mssing image name")
		}
		containerId := ctx.Args().Get(0)
		imageName := ctx.Args().Get(1)
		return commitContainer(containerId, imageName)
	},
}

func commitContainer(containerId, imageName string) error {
	return container.Commit(containerId, imageName)
}
