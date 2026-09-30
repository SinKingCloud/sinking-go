package daemon

import (
	"fmt"
	"os"
)

// readPidFile 读取守护进程 PID。
func (u *Daemon) readPidFile() (int, error) {
	data, err := os.ReadFile(u.PidFileName)
	if err != nil {
		return 0, err
	}
	var pid int
	if _, err = fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return 0, err
	}
	return pid, nil
}
