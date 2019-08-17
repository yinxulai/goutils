package terminal

import (
	"io"
	"os/exec"
)

// New 创建一个绘画
func New() (*Session, error) {
	var err error
	session := new(Session)
	session.session = exec.Command("echo", "start....")
	session.stdin, err = session.session.StdinPipe()
	session.stdout, err = session.session.StdoutPipe()
	session.stderr, err = session.session.StderrPipe()
	err = session.session.Start()
	if err != nil {
		return nil, err
	}
	return session, err
}

// Session Session
type Session struct {
	session *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
}

func (s *Session) Env(stdout io.ReadCloser, stderr io.ReadCloser, stdin io.WriteCloser) {

}

// Action 一个动作
// stdout 上个命令的输出
// stderr 上个命令的错误输出
// stdin 输入
func (s *Session) Action(stdout io.ReadCloser, stderr io.ReadCloser, stdin io.WriteCloser) {

}
