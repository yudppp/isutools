package profile

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/yudppp/isutools/utils/slackcat"
)

var once sync.Once

// StartAll .
func StartAll(duration time.Duration) error {
	err := StartMem(duration)
	if err != nil {
		return err
	}
	return StartCPU(duration, false)
}

// StartCPU .
func StartCPU(duration time.Duration, imageOnly bool) error {
	f, err := ioutil.TempFile("", "cpu")
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}
	timerStop := func() {
		time.Sleep(duration)
		pprof.StopCPUProfile()
		f.Close()

		// image format
		if true {
			imageFile := "pprof.png"
			exec.Command("go", "tool", "pprof", "-png", "-output", imageFile, f.Name()).Run()
			slackcat.SendFile(imageFile, "pprof.png")
			os.Remove(imageFile)
		} else {
			imageFile := "pprof.svg"
			exec.Command("go", "tool", "pprof", "-svg", "-output", imageFile, f.Name()).Run()
			slackcat.SendFile(imageFile, "pprof.svg")
			os.Remove(imageFile)
		}

		if !imageOnly {
			slackcat.SendFile(f.Name(), "cpu.pprof")
		}
		os.Remove(f.Name())
	}
	go timerStop()
	return nil
}

// StartMem .
func StartMem(duration time.Duration) error {
	f, err := ioutil.TempFile("", "mem")
	if err != nil {
		return err
	}
	timerStop := func() {
		time.Sleep(duration)
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
		slackcat.SendFile(f.Name(), "mem.mprof")
		os.Remove(f.Name())
	}
	go timerStop()
	return nil
}
