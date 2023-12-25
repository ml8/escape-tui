package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ml8/escape-tui/model"
	esc "github.com/ml8/escape-tui/model"

	c "github.com/fatih/color"
	"github.com/golang/glog"
)

var (
	stateFile = flag.String("states", "states.yaml", "state config file")
)

func main() {
	flag.Parse()

	glog.Infof("Using %v for states...", stateFile)

	s, err := ioutil.ReadFile(*stateFile)
	if err != nil {
		glog.Fatalf("Error reading %v: %v", *stateFile, err)
	}

	fmt.Print("\033[H\033[2J")

	out := model.TypewriteWith(c.New(c.FgHiGreen).Add(c.Bold).Printf)
	aside := model.TypewriteWith(c.New(c.FgMagenta).Printf)
	er := model.TypewriteWith(c.New(c.FgHiRed).Add(c.Bold).Printf)

	model := esc.Parse(esc.StdIn(), esc.OutFrom(out, aside, er), string(s))
	model.Run()

	c.Unset()
}
