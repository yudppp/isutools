package slackcat

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"time"
)

var defaultChannel = "notification"

// SendFile .
func SendFile(filePath, pattern string) error {
	ext := path.Ext(pattern)
	names := strings.Split(pattern, ".")
	outputName := fmt.Sprintf("%s-%s%s", names[0], time.Now().Format("200601021504"), ext)
	return exec.Command("slackcat", "--channel", "notification", "-n", outputName, filePath).Run()
}

// SetDefaultChannel .
func SetDefaultChannel(channel string) {
	defaultChannel = channel
}
