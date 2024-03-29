package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	// 需要导入nsenter包，以触发C代码
	"github.com/pjimming/mydocker/command"
	_ "github.com/pjimming/mydocker/nsenter"
)

const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		command.InitCommand,
		command.CommitCommand,
		command.ExecCommand,
		command.ListCommand,
		command.RemoveCommand,
		command.LogCommand,
		command.RunCommand,
		command.StopCommand,
		command.NetworkCommand,
	}

	app.Before = func(ctx *cli.Context) error {
		// init logger settings
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
		})
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
