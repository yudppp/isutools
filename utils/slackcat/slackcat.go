package slackcat

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	// "bytes"
	"os/exec"
)

var defaultChannel = "notification"
var tokenOption = ""

// SendFile .
func SendFile(filePath, pattern string) error {
	ext := path.Ext(pattern)
	names := strings.Split(pattern, ".")
	outputName := fmt.Sprintf("%s-%s%s", names[0], time.Now().Format("200601021504"), ext)
	fmt.Println("slackcat", "--channel", defaultChannel, tokenOption, "-n", outputName, filePath)
	return exec.Command("slackcat", "--channel", defaultChannel, tokenOption, "-n", outputName, filePath).Run()
}

// SendText .
func SendText(filename, text string) error {
	c1 := exec.Command("echo", text)
	c2 := exec.Command("slackcat", "--channel", defaultChannel, tokenOption, "--filename", filename)

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
	tokenOption = fmt.Sprintf("--token %s", token)
}
