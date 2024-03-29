package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/pjimming/mydocker/container"
)

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "list all of the containers",
	Action: func(ctx *cli.Context) error {
		listContainers()
		return nil
	},
}

// listContainers 获取所有容器信息，并且打印出来
// 首先遍历存放容器数据的/var/lib/mydocker/containers/目录，里面每一个子目录都是一个容器。
// 然后使用 getContainerInfo 方法解析子目录中的 config.json 文件拿到容器信息
// 最后格式化成 table 形式打印出来即可
func listContainers() {
	dirs, err := os.ReadDir(container.InfoLoc)
	if err != nil {
		logrus.Errorf("[listContainers] read dir fail, %v", err)
		return
	}

	containers := make([]*container.Info, 0)
	for _, dir := range dirs {
		containerInfo, err := container.ReadInfo(dir.Name())
		if err != nil {
			logrus.Errorf("[listContainers] read %s info fail, %v", dir.Name(), err)
			continue
		}
		containers = append(containers, containerInfo)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	if _, err = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n"); err != nil {
		logrus.Errorf("[listContainers] Fprint fail, %v", err)
	}

	for _, item := range containers {
		if _, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime,
		); err != nil {
			logrus.Errorf("[listContainers] Fprint fail %v", err)
		}
	}
	if err = w.Flush(); err != nil {
		logrus.Errorf("[listContainers] tabwriter flush error, %v", err)
	}
}
