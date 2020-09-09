package slackcat

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"os/exec"
)

var defaultChannel = "notification"
var tokenOption = ""

// SendFile .
func SendFile(filePath, pattern string) error {
	ext := path.Ext(pattern)
	names := strings.Split(pattern, ".")
	outputName := fmt.Sprintf("%s-%s%s", names[0], time.Now().Format("200601021504"), ext)
	if tokenOption == "" {
		return exec.Command("slackcat", "--channel", defaultChannel, "-n", outputName, filePath).Run()
	}
	return exec.Command("slackcat", "--channel", defaultChannel, "--token", tokenOption, "-n", outputName, filePath).Run()
}

// SendText .
func SendText(filename, text string) error {
	c1 := exec.Command("echo", text)
	var command2Options []string
	if tokenOption == "" {
		command2Options = []string{"--channel", defaultChannel, "--token", tokenOption, "--filename", filename}
	} else {
		command2Options = []string{"--channel", defaultChannel, "--filename", filename}
	}
	c2 := exec.Command("slackcat", command2Options...)

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	return c2.Wait()
}

// SetDefaultChannel .
func SetDefaultChannel(channel string) {
	defaultChannel = channel
}

// SetToken .
func SetToken(token string) {
	if token == "" {
		tokenOption = ""
	}
	tokenOption = token
}
