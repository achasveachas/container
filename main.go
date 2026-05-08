package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Entry point: dispatches to run or child mode based on the first argument.
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Invalid argument")
	}
}

// run launches a new child process in a new set of namespaces.
// It re-executes itself with "child" as the first argument, passing through any additional arguments.
func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	// Prepare to re-execute this binary with "child" as the first argument
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Connect standard input/output/error to the new process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up Linux namespaces: UTS (hostname), PID, and mount
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// Start the child process
	must(cmd.Run())
}


// child runs inside the new namespaces and chroot, acting as the isolated container process.
func child() {
	fmt.Printf("Running from child process %v\n", os.Args[2:])

	// Set a new hostname (UTS namespace)
	must(syscall.Sethostname([]byte("container")))

	// Change root to the ubuntu-rootfs directory (filesystem isolation)
	pwd, err := os.Getwd()
	must(err)

	must(syscall.Chroot(filepath.Join(pwd, "ubuntu-rootfs")))
	must(os.Chdir("/"))

	// Mount /proc for process info inside the container
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("proc", 0)
	
	// Execute the requested command inside the container
	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}


// must is a helper that panics if an error occurs, simplifying error handling.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
