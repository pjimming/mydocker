package command

import (
	"fmt"
	"github.com/pjimming/mydocker/container"
	"github.com/urfave/cli"
)

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: "Commit container to image",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("mssing image name")
		}
		imageName := ctx.Args().Get(0)
		return commitContainer(imageName)
	},
}

func commitContainer(imageName string) error {
	return container.Commit(imageName)
}
