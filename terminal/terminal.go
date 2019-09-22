package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// New 创建一个绘画
func New() (*Terminal, error) {
	var err error
	terminal := new(Terminal)
	terminal.Terminal = exec.Command("bash")
	terminal.inWriteCloser, err = terminal.Terminal.StdinPipe() // 输入
	terminal.outReadCloser, err = terminal.Terminal.StdoutPipe()
	terminal.errReadCloser, err = terminal.Terminal.StderrPipe()
	terminal.outRender = bufio.NewReader(terminal.outReadCloser)
	terminal.errRender = bufio.NewReader(terminal.errReadCloser)
	return terminal, err
}

// Terminal Terminal
type Terminal struct {
	runing        bool
	Terminal      *exec.Cmd
	outRender     io.Reader
	errRender     io.Reader
	inWriteCloser io.WriteCloser
	outReadCloser io.ReadCloser
	errReadCloser io.ReadCloser
}

func (s *Terminal) start() error {
	if !s.runing {
		if err := s.Terminal.Start(); err != nil {
			return err
		}
		s.runing = true
	}
	return nil
}

// GetEnv 获取环境变量
// 仅第一个 Action 执行前可用
func (s *Terminal) GetEnv(key string) string {
	for _, env := range s.Terminal.Env {
		if key != "" && env[:len(key)] == key {
			return env[len(key)+1:]
		}
	}
	return ""
}

// SetEnv 设置环境变量
// 仅第一个 Action 执行前可用
func (s *Terminal) SetEnv(key, value string) {
	s.Terminal.Env = append(s.Terminal.Env, fmt.Sprintf("%s=%s", key, value))
}

// InputString 输入字符串
func (s *Terminal) InputString(input string, timeout time.Duration) (int, string, error) {
	return s.InputBytes([]byte(input), timeout)
}

// InputBytes 输入字节
func (s *Terminal) InputBytes(input []byte, timeout time.Duration) (int, string, error) {
	var err error
	inputBuf := new(bytes.Buffer)

	err = s.start()
	if err != nil {
		return 0, "", err
	}

	inputBuf.Write(input)
	inputBuf.WriteString("\n") // 自动添加回车？
	_, err = s.inWriteCloser.Write(inputBuf.Bytes())
	if err != nil {
		return 0, "", err
	}

	for {
		// outdata, _, outerr := s.outRender.Read()
		// if outerr != nil || outerr == io.EOF {
		// 	// 发生错或者读到 EOF
		// 	break
		// }
		// fmt.Printf("out: %s,%s \n", outputBuffer.Bytes(), errputBuffer.Bytes())
	}

	fmt.Printf("end")
	return 0, "", nil
}

// InputRune 输入 Rune
func (s *Terminal) InputRune(input []rune, timeout time.Duration) (int, string, error) {
	return s.InputString(string(input), timeout)
}

// Exit 退出
func (s *Terminal) Exit(code uint) error {
	return nil
}
