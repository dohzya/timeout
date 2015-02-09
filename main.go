package main

/* TODO
 *
 * Accept floating numbers
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

func parseDuration(s string) (*time.Duration, error) {
	dur := s[:len(s)-1]
	var mod time.Duration
	c := s[len(s)-1]
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		mod = time.Second
		dur = s
	case 's':
		mod = time.Second
	case 'm':
		mod = time.Minute
	case 'h':
		mod = time.Hour
	case 'd':
		mod = 24 * time.Hour
	default:
		return nil, fmt.Errorf("Bad modifier: %v (expected: s|m|h|d)", c)
	}
	timeout, err := strconv.ParseInt(dur, 0, 64)
	if err != nil {
		return nil, err
	}
	duration := time.Duration(timeout) * mod
	return &duration, nil
}

func main() {
	signame := flag.String("s", "15", "The signal to use")
	flag.Parse()

	if len(flag.Args()) < 1 {
		die("[timeout] Missing timeout")
	}
	timeout, timeouterr := parseDuration(flag.Arg(0))
	if timeouterr != nil {
		die("[timeout] Bad timeout value")
	}

	cmd := flag.Args()[1:]
	if len(cmd) == 0 {
		die("[timeout] Missing command")
	}

	sig, sigerr := getSignal(*signame)
	if sigerr {
		die(fmt.Sprintf("[timeout] Unknown signal: %v", signame))
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		die2(fmt.Sprintf("[timeout] Can't start the process: %v", err), 127)
	}

	timer := time.AfterFunc(*timeout, func() {
		if err := command.Process.Signal(sig); err != nil {
			fmt.Fprintf(os.Stderr, "[timeout] Can't kill the process: %v\n", err.Error())
		}
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
