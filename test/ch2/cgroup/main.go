package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const (
	// 挂载memory subsystem的hierarchy的根目录位置
	cgroupMemoryHierarchyMount = "/sys/fs/cgroups/memory"
	testMemoryLimit            = "testmemorylimit"
)

func main() {
	// 被克隆时，0会被设置为/proc/self/exe
	if os.Args[0] == "/proc/self/exe" {
		log.Printf("current pid %d", syscall.Getpid())
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	// 克隆一个与自己相同的一个进程
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 异步启动，不会阻塞
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		// 得到fork出来进程映射在外部命名空间的pid
		log.Printf("%d", cmd.Process.Pid)
		// 在系统默认创建挂载了memory subsystem的hierarchy上创建cgroup
		_ = os.Mkdir(path.Join(cgroupMemoryHierarchyMount, testMemoryLimit), 0755)
		// 将当前进程添加到cgroup
		_ = os.WriteFile(path.Join(cgroupMemoryHierarchyMount, testMemoryLimit, "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// 限制cgroup进程使用
		_ = os.WriteFile(path.Join(cgroupMemoryHierarchyMount, testMemoryLimit, "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
	_, _ = cmd.Process.Wait()
}
