package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/atpons/octool/cmd/octool/config"

	"os/exec"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	config.Load(os.Getenv("CONFIG_FILE"))
	ocproxy := NewTmpFile("ocproxy.sh")
	go func() {
		p := readPassword()

		_ = createOcproxyScript(ocproxy.file)

		c := exec.Command("openconnect", config.BuildOpenConnectOpts(ocproxy.file.Name())...)
		c.Stdin = strings.NewReader(p)

		cSt := exec.Command("straightforward", config.BuildStraightForwardOpts()...)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		if config.Value.StraightForward.Enabled {
			go runOpenConnect(c, cSt)
		} else {
			go run(c)
		}

		wg.Wait()
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)
	<-sig
	fmt.Println("closing")
	ocproxy.Close()
}

func typePassword() string {
	var password []byte
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	return string(password)
}

func readPassword() string {
	var password string
	if terminal.IsTerminal(syscall.Stdin) {
		password = typePassword()
	} else {
		p, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		password = string(p)
	}
	return password
}

func run(cmd *exec.Cmd) {
	stderr, _ := cmd.StdoutPipe()
	stdout, _ := cmd.StderrPipe()

	cmd.Start()
	go showOutput(stdout)
	go showOutput(stderr)
	cmd.Wait()
}

func runOpenConnect(cmd, afterCmd *exec.Cmd) {
	stderr, _ := cmd.StdoutPipe()
	stdout, _ := cmd.StderrPipe()

	cmd.Start()
	go hookOpenConnect(stdout, afterCmd)
	go hookOpenConnect(stderr, afterCmd)
	cmd.Wait()
}

func showOutput(p io.ReadCloser) {
	scanner := bufio.NewScanner(p)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func hookOpenConnect(p io.ReadCloser, cmd *exec.Cmd) {
	scanner := bufio.NewScanner(p)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		// "Connected as" describes establishing connection with VPN server, and go into other command.
		if strings.Contains(m, "Connected as") {
			go run(cmd)
		}
	}
}

type tmpFile struct {
	file *os.File
}

func (t *tmpFile) Close() {
	_ = os.Remove(t.file.Name())
}

func NewTmpFile(f string) *tmpFile {
	t := tmpFile{}
	t.file, _ = ioutil.TempFile("", f)

	return &t
}

func createOcproxyScript(f *os.File) error {
	defer f.Close()

	output := "#!/bin/sh\nocproxy -D " + config.Value.OcProxy.Port
	_, err := f.Write(([]byte)(output))
	if err != nil {
		return err
	}

	err = os.Chmod(f.Name(), 0755)
	if err != nil {
		return err
	}

	return nil
}
