package ps

import (
	"fmt"
	"os/exec"
)

func ExecutePsCommandByPid(pid int) (string, error) {
	// The -o comm= option tells ps to only output the command path of the process, without any headers or additional formatting
	// The -p option filters the output to only include the process with the specified PID
	prompt := fmt.Sprintf("ps -p %d -o comm=", pid)
	cmd := exec.Command("sh", "-c", prompt)

	output, error := cmd.Output()
	if error != nil {
		return "", error
	}

	return string(output), nil
}

func KillProcessByPid(pid int) error {
	// The -9 option sends the SIGKILL signal to forcefully terminate the process
	prompt := fmt.Sprintf("kill -9 %d", pid)
	cmd := exec.Command("sh", "-c", prompt)

	_, error := cmd.Output()
	if error != nil {
		return error
	}

	return nil
}
