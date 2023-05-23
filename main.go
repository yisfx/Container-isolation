package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("arg undefined")
	}
}

func run() {

	cmd := exec.Command(os.Args[0], append([]string{"child"}, os.Args[2:]...)...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			// syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			// syscall.CLONE_NEWUSER |
			// syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,

		// UidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 1,
		// 		HostID:      0,
		// 		Size:        1,
		// 	},
		// },
		// GidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 1,
		// 		HostID:      0,
		// 		Size:        1,
		// 	},
		// },
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("cant run", err)
	}
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	var cmd *exec.Cmd
	{
		// cmd = exec.Command("pwd")
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Run()

		// cmd = exec.Command("ls")
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Run()
		// return
	}

	defer func() {
		syscall.Unmount("/proc", syscall.MNT_FORCE)
		fmt.Println("end!!!!!")
	}()
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))
	must(syscall.Sethostname([]byte("mycontainer")))
	// must(os.Mkdir("rootfs/oldrootfs", 0700))
	syscall.PivotRoot("rootfs", "rootfs/oldrootfs")
	must(syscall.Chroot("./rootfs"))
	must(syscall.Chdir("/"))

	cmd = exec.Command(os.Args[2])
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	fmt.Println("Start Run")
	must(cmd.Run())
}
func must(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
