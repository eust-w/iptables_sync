package ctrl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"text/template"
)

type Bash struct {
	Command   string
	PipeFail  bool
	Arguments map[string]string

	retCode int
	stdout  string
	stderr  string
	err     error
}

func (b *Bash) build() error {
	if b.Arguments != nil {
		tmpl, err := template.New("script").Parse(b.Command)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, b.Arguments)
		if err != nil {
			return err
		}

		b.Command = buf.String()
	}

	if b.PipeFail {
		b.Command = fmt.Sprintf("set -o pipefail; %s", b.Command)
	}

	return nil
}

func (b *Bash) RunWithReturn() (retCode int, stdout, stderr string, err error) {
	if err = b.build(); err != nil {
		b.err = err
		return -1, "", "", err
	}

	var so, se bytes.Buffer
	var cmd *exec.Cmd

	if len(b.Command) > 1024*4 {
		content := []byte(b.Command)
		tmpfile, _ := ioutil.TempFile("/tmp", "zsnagent")
		err = tmpfile.Chmod(0777)
		_, err = tmpfile.Write(content)
		err = tmpfile.Close()
		cmd = exec.Command("bash", "-c", tmpfile.Name())
		defer os.Remove(tmpfile.Name())
	} else {
		cmd = exec.Command("bash", "-c", b.Command)
	}

	cmd.Stdout = &so
	cmd.Stderr = &se

	var waitStatus syscall.WaitStatus
	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			retCode = waitStatus.ExitStatus()
		} else {
			err = fmt.Errorf("unable to get return code, %s: %s", err, so.String()+se.String())
			return
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		retCode = waitStatus.ExitStatus()
	}

	stdout = so.String()
	stderr = se.String()

	b.retCode = retCode
	b.stdout = stdout
	b.stderr = stderr

	return
}

func NewBash() *Bash {
	return &Bash{}
}
