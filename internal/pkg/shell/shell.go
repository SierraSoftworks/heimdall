package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Shell struct {
	Command     string
	Args        []string
	Directory   string
	Environment []string
}

func NewDefaultShell() *Shell {
	switch runtime.GOOS {
	case "windows":
		return NewPowerShell()
	case "linux":
		return NewBashShell()
	default:
		return NewShShell()
	}
}

func NewShShell() *Shell {
	return &Shell{
		Command: "/bin/sh",
		Args: []string{
			"-c",
		},
		Directory:   "/tmp",
		Environment: []string{},
	}
}

func NewBashShell() *Shell {
	return &Shell{
		Command: "/bin/bash",
		Args: []string{
			"-c",
		},
		Directory:   "/tmp",
		Environment: []string{},
	}
}

func NewPowerShell() *Shell {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = "C:\\Windows"
	}

	return &Shell{
		Command: "powershell.exe",
		Args: []string{
			"-ExecutionPolicy", "Unrestricted",
			"-NonInteractive",
			"-NoProfile",
			"-Command",
		},
		Directory: os.Getenv("TEMP"),
		Environment: []string{
			// Fix PowerShell Error 8009001d when SystemRoot env variable is unset
			fmt.Sprintf("SystemRoot=%s", systemRoot),
		},
	}
}

func (s *Shell) Describe() string {
	return s.Command
}

func (s *Shell) NewCommand(command string, env []string) (*ShellCommand, error) {
	cmd := exec.Command(s.Command, append(s.Args, command)...)
	cmd.Dir = os.ExpandEnv(s.Directory)
	cmd.Env = append(s.Environment, env...)

	sc := &ShellCommand{
		cmd,
	}

	sc.configure()

	return sc, nil
}
