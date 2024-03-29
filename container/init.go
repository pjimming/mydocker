package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	readPipeFdIndex = 3
)

// RunContainerInitProcess 启动容器的init进程
/*
这里的init函数是在容器内部执行的，也就是说，代码执行到这里后，容器所在的进程其实就已经创建出来了，
这是本容器执行的第一个进程。
使用mount先去挂载proc文件系统，以便后面通过ps等系统命令去查看当前进程资源的情况。
*/
func RunContainerInitProcess() error {
	// mount -t proc proc /proc
	mountProc()

	// read pipe
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) <= 0 {
		return fmt.Errorf("run container get user command fail, command array is nil")
	}

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)

	if err = syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf("exec command fail, %v", err)
		return err
	}
	return nil
}

func readUserCommand() []string {
	// uintptr(3)就是指 index 为3的文件描述符，也就是传递进来的管道的另一端，至于为什么是3，具体解释如下：
	/*	因为每个进程默认都会有3个文件描述符，分别是标准输入、标准输出、标准错误。这3个是子进程一创建的时候就会默认带着的，
		前面通过ExtraFiles方式带过来的 readPipe 理所当然地就成为了第4个。
		在进程中可以通过index方式读取对应的文件，比如
		index0：标准输入
		index1：标准输出
		index2：标准错误
		index3：带过来的第一个FD，也就是readPipe
		由于可以带多个FD过来，所以这里的3就不是固定的了。
		比如像这样：cmd.ExtraFiles = []*os.File{a,b,c,readPipe} 这里带了4个文件过来，分别的index就是3,4,5,6
		那么我们的 readPipe 就是 index6,读取时就要像这样：pipe := os.NewFile(uintptr(6), "pipe")
	*/
	pipe := os.NewFile(uintptr(readPipeFdIndex), "pipe")
	defer func() {
		_ = pipe.Close()
	}()

	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("read pipe fail, %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

func mountProc() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("os getwd fail, %v", err)
		return
	}

	logrus.Infof("Current location is %s", pwd)

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示声明你要这个新的mount namespace独立。
	// 即 mount proc 之前先把所有挂载点的传播类型改为 private，避免本 namespace 中的挂载事件外泄。
	// 把所有挂载点的传播类型改为 private，避免本 namespace 中的挂载事件外泄。
	_ = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	if err = pivotRoot(pwd); err != nil {
		logrus.Errorf("pivot_root fail, %v", err)
		return
	}

	// 如果不先做 private mount，会导致挂载事件外泄，后续再执行 mydocker 命令时 /proc 文件系统异常
	// 可以执行 mount -t proc proc /proc 命令重新挂载来解决

	// MS_NOEXEC 在本文件系统不允许运行其他程序。
	// MS_NOSUID 在本系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID
	// MS_NOD 这个参数是自 Linux 2.4 ，所有 mount 的系统都会默认设定的参数。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
}

func pivotRoot(root string) error {
	// NOTE：PivotRoot调用有限制，newRoot和oldRoot不能在同一个文件系统下。
	// 因此，为了使当前root的老root和新root不在同一个文件系统下，这里把root重新mount了一次。
	// bind mount是把相同的内容换了一个挂载点的挂载方法
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		logrus.Errorf("mount rootfs to itself fail, %v", err)
		return err
	}
	// 创建rootfs/.pivot_root存储old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", pivotDir, err)
		return err
	}
	// 执行pivot_root调用,将系统rootfs切换到新的rootfs,
	// PivotRoot调用会把 old_root挂载到pivotDir,也就是rootfs/.pivot_root,挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		logrus.Errorf("pivot_root fail, %v", err)
		return err
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		logrus.Errorf("chdir fail, %v", err)
		return err
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// 最后再把old_root umount了，即 umount rootfs/.pivot_root
	// 由于当前已经是在 rootfs 下了，就不能再用上面的rootfs/.pivot_root这个路径了,现在直接用/.pivot_root这个路径即可
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		logrus.Errorf("unmount pivot_root dir fail, %v", err)
		return err
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
