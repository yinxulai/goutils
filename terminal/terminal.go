package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

// ActionFunc ActionFunc
type ActionFunc = func([]byte, []byte, io.WriteCloser)

// New 创建一个绘画
func New() (*Session, error) {
	var err error
	session := new(Session)
	session.session = exec.Command("bash", "-c")
	session.stdin, err = session.session.StdinPipe()
	session.stdout, err = session.session.StdoutPipe()
	session.stderr, err = session.session.StderrPipe()
	session.actions = make(map[*regexp.Regexp]ActionFunc, 10)
	return session, err
}

// Session Session
type Session struct {
	runing  bool
	session *exec.Cmd
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	stdin   io.WriteCloser
	actions map[*regexp.Regexp]ActionFunc
}

// GetEnv 获取环境变量
// 仅第一个 Action 执行前可用
func (s *Session) GetEnv(key string) string {
	for _, env := range s.session.Env {
		if key != "" && env[:len(key)] == key {
			return env[len(key)+1:]
		}
	}
	return ""
}

// SetEnv 设置环境变量
// 仅第一个 Action 执行前可用
func (s *Session) SetEnv(key, value string) {
	s.session.Env = append(s.session.Env, fmt.Sprintf("%s=%s", key, value))
}

// action 一个动作
// stdout 上个命令的输出
// stderr 上个命令的错误输出
// stdin 输入
func (s *Session) action(handler func(io.ReadCloser, io.ReadCloser, io.WriteCloser)) {
	if !s.runing {
		if err := s.session.Run(); err != nil {
			s.Exit(1)
		}
	}
	handler(s.stdout, s.stderr, s.stdin)
}

// AddAction 注册匹配自执行动作
func (s *Session) AddAction(reg regexp.Regexp, handler ActionFunc) {
	s.actions[&reg] = handler
}

// Run 退出
func (s *Session) Run(code uint) error {
	if !s.runing {
		if err := s.session.Start(); err != nil {
			s.Exit(1)
			return err
		}

		outReader := bufio.NewReader(s.stdout)
		errReader := bufio.NewReader(s.stderr)

		// 读取内容
		for {
			outdata, _, outerr := outReader.ReadRune()
			if outerr != nil || outerr == io.EOF {
				// 发生错或者读到 EOF
				break
			}
			errdata, _, errerr := errReader.ReadRune()
			if errerr != nil || errerr == io.EOF {
				// 发生错或者读到 EOF
				break
			}

			for reg, handler := range s.actions {
				if reg.MatchString(string(errdata)) || reg.MatchString(string(outdata)) {
					handler([]byte(string(errdata)), []byte(string(outdata)), s.stdin)
				}
			}

		}
	}

	return nil
}

// Exit 退出
func (s *Session) Exit(code uint) error {
	_, err := s.stdin.Write([]byte(fmt.Sprintf("\n exit %d \n", code)))
	return err
}
