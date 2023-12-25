package model

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type PrintFunc func(string, ...any) (int, error)

var (
	rng = rand.New(rand.NewSource(time.Now().Unix()))
	min = 10000
	max = 100000000
)

type In interface {
	ReadString(delim byte) (s string, err error)
}

type Out interface {
	WriteOut(format string, a ...any) (n int, err error)
	WriteAside(format string, a ...any) (n int, err error)
	WriteErr(format string, a ...any) (n int, err error)
}

type owrapper struct {
	out   PrintFunc
	aside PrintFunc
	err   PrintFunc
}

func TypewriteWith(out PrintFunc) PrintFunc {
	return func(format string, a ...any) (n int, err error) {
		data := fmt.Sprintf(format, a...)
		for _, c := range data {
			n, err = out(string(c))
			time.Sleep(time.Duration(min) + time.Duration(rng.Int63n(int64(max-min))))
		}
		return n, err
	}
}

func (o owrapper) WriteOut(format string, a ...any) (n int, err error) {
	return o.out(format, a...)
}

func (o owrapper) WriteAside(format string, a ...any) (n int, err error) {
	return o.aside(format, a...)
}

func (o owrapper) WriteErr(format string, a ...any) (n int, err error) {
	return o.err(format, a...)
}

func OutFrom(o, a, e PrintFunc) Out {
	return owrapper{o, a, e}
}

type stdin struct {
	reader *bufio.Reader
}
type stdout struct{}

func (t stdin) ReadString(delim byte) (s string, err error) {
	s, err = t.reader.ReadString(delim)
	s = strings.TrimSpace(s)
	return
}

func (s stdout) WriteOut(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}
func (s stdout) WriteAside(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}
func (s stdout) WriteErr(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}
func StdIn() In {
	return stdin{bufio.NewReader(os.Stdin)}
}
func StdOut() Out {
	return stdout{}
}
