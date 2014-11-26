Timeout
=======

Re-implementation of the Linux's “timeout” command

Runs the specified command but kill-it after the specified timeout.

By default the SIGTERM is sent to the process.

Use
---

#### Basic version

```bash
timeout 3 /my/command/to/run arg1 arg2
```

#### With custom time

Here: 3 minutes

```bash
timeout 3m /my/command/to/run arg1 arg2
```

Accepted modifiers: “s”, “m”, “h”, “d”

#### Specifying the signal to use

Here: SIGKILL

```bash
timeout -s 9 3 /my/command/to/run arg1 arg2
```

Accepted signals:

- “1” or “HUP”
- “2” or “INT”
- “3” or “QUIT”
- “4” or “ILL”
- “5” or “TRAP”
- “6” or “ABRT”
- “7” or “BUS”
- “8” or “FPE”
- “9” or “KILL”
- “10” or “USR1”
- “11” or “SEGV”
- “12” or “USR2”
- “13” or “PIPE”
- “14” or “ALRM”
- “15” or “TERM”
- “17” or “CHLD”
- “18” or “CONT”
- “19” or “STOP”
- “20” or “TSTP”
- “21” or “TTIN”
- “22” or “TTOU”
- “23” or “URG”
- “24” or “XCPU”
- “25” or “XFSZ”
- “26” or “VTALRM”
- “27” or “PROF”
- “28” or “WINCH”
- “29” or “IO”
- “31” or “SYS”
