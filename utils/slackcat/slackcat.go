package slackcat

import (
	"fmt"
	"path"
	"strings"
	"time"
	"io"
	// "bytes"
	"os/exec"
)

var defaultChannel = "notification"

// SendFile .
func SendFile(filePath, pattern string) error {
	ext := path.Ext(pattern)
	names := strings.Split(pattern, ".")
	outputName := fmt.Sprintf("%s-%s%s", names[0], time.Now().Format("200601021504"), ext)
	return exec.Command("slackcat", "--channel", defaultChannel, "-n", outputName, filePath).Run()
}

// SendText .
func SendText(filename, text string) error {
	c1 := exec.Command("echo", text)
    c2 := exec.Command("slackcat", "--channel", defaultChannel, "--filename", filename)

    r, w := io.Pipe()
    c1.Stdout = w
	c2.Stdin = r
	
	// var out bytes.Buffer
    // c2.Stdout = &out

    c1.Start()
    c2.Start()
    c1.Wait()
    w.Close()
    return  c2.Wait()		
	// return nil
}

// SetDefaultChannel .
func SetDefaultChannel(channel string) {
	defaultChannel = channel
}
