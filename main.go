package main

/* TODO
 *
 * DURATION is a floating point number with an optional suffix: 's' for seconds
 * (the default), 'm' for minutes, 'h' for hours or 'd' for days.
 */

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func die(msg string) {
	die2(msg, 125)
}

func die2(msg string, status int) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(status)
}

func getSignal(name string) (syscall.Signal, bool) {
	switch name {
	default:
		return syscall.SIGTERM, true
	case "1", "HUP":
		return syscall.SIGHUP, false
	case "2", "INT":
		return syscall.SIGINT, false
	case "3", "QUIT":
		return syscall.SIGQUIT, false
	case "4", "ILL":
		return syscall.SIGILL, false
	case "5", "TRAP":
		return syscall.SIGTRAP, false
	case "6", "ABRT":
		return syscall.SIGABRT, false
	case "7", "BUS":
		return syscall.SIGBUS, false
	case "8", "FPE":
		return syscall.SIGFPE, false
	case "9", "KILL":
		return syscall.SIGKILL, false
	case "10", "USR1":
		return syscall.SIGUSR1, false
	case "11", "SEGV":
		return syscall.SIGSEGV, false
	case "12", "USR2":
		return syscall.SIGUSR2, false
	case "13", "PIPE":
		return syscall.SIGPIPE, false
	case "14", "ALRM":
		return syscall.SIGALRM, false
	case "15", "TERM":
		return syscall.SIGTERM, false
	case "17", "CHLD":
		return syscall.SIGCHLD, false
	case "18", "CONT":
		return syscall.SIGCONT, false
	case "19", "STOP":
		return syscall.SIGSTOP, false
	case "20", "TSTP":
		return syscall.SIGTSTP, false
	case "21", "TTIN":
		return syscall.SIGTTIN, false
	case "22", "TTOU":
		return syscall.SIGTTOU, false
	case "23", "URG":
		return syscall.SIGURG, false
	case "24", "XCPU":
		return syscall.SIGXCPU, false
	case "25", "XFSZ":
		return syscall.SIGXFSZ, false
	case "26", "VTALRM":
		return syscall.SIGVTALRM, false
	case "27", "PROF":
		return syscall.SIGPROF, false
	case "28", "WINCH":
		return syscall.SIGWINCH, false
	case "29", "IO":
		return syscall.SIGIO, false
	case "31", "SYS":
		return syscall.SIGSYS, false
	}
}

func main() {
	signame := flag.String("s", "15", "The signal to use")
	flag.Parse()

	timeout := (func() int64 {
		if len(flag.Args()) < 1 {
			die("[timeout] Missing timeout")
		}
		arg0 := flag.Arg(0)
		timeout, err := strconv.ParseInt(arg0, 0, 64)
		if err != nil {
			die("[timeout] Bad timeout value")
		}
		return timeout
	})()

	cmd := flag.Args()[1:]
	if len(cmd) == 0 {
		die("[timeout] Missing command")
	}

	sig, sigerr := getSignal(*signame)
	if sigerr {
		die(fmt.Sprintf("[timeout] Unknown signal: %v", signame))
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		die2(fmt.Sprintf("[timeout] Can't start the process: %v", err), 127)
	}

	timer := time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		command.Process.Signal(syscall.SIGTERM)
	})

	err = command.Wait()

	killed := !timer.Stop()

	status := 0
	if killed {
		if sig == syscall.SIGKILL {
			status = 132
		} else {
			status = 124
		}
	} else if err != nil {
		if command.ProcessState == nil {
			die2(fmt.Sprintf("[timeout] Error occured: %v", err), 127)
		}
		status = command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	} else {
		status = 0
	}
	os.Exit(status)

}
