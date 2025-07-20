package lib

import (
	"os/exec"
	"runtime"
	"strings"
)

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	case "darwin":
		cmd = "open"
		args = []string{}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start"}
		} else {
			cmd = "xdg-open"
			args = []string{}
		}
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
