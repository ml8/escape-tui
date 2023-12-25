package model

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Tag string

type gamemodel struct {
	States []*AnswerSet // Set of answers
	Tags   []Tag        // Initial tags
}

type AnswerSet struct {
	Accept   map[string]string `yaml:",omitempty"`
	Partial  map[string]string `yaml:",omitempty"`
	Strict   bool              `yaml:",omitempty"`
	Final    bool              `yaml:",omitempty"`
	Webhook  string            `yaml:",omitempty"`
	Requires []Tag             `yaml:",omitempty"`
	Provides []Tag             `yaml:",omitempty"`
	Consumes []Tag             `yaml:",omitempty"`
}

// XXX todo and webhook

type ResultType int

const (
	Success = iota
	PartialSuccess
	Failure
)

type Result struct {
	Type  ResultType
	Final bool
	Txt   string
	Url   string
}

type PlayerState struct {
	tags []Tag
}

type Game struct {
	i         In
	o         Out
	answerSet []*AnswerSet
	state     *PlayerState
}

func Parse(in In, out Out, model string) *Game {
	parsed := parse(model)
	return &Game{i: in, o: out, answerSet: parsed.States, state: &PlayerState{tags: parsed.Tags}}
}

func webhook(url string) {
	// TODO retry
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Fatalf("Could not create request: %s\n", err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		glog.Fatalf("Error making http request: %s\n", err)
	}
}

func (g *Game) TryAll(str string) (r *Result) {
	for _, s := range g.answerSet {
		r = s.Try(str, g.state)
		if r.Type != Failure {
			return r
		}
	}
	return r
}

func (g *Game) Run() {
	var s string
	for {
		g.o.WriteOut("> ")
		s, _ = g.i.ReadString('\n')
		r := g.TryAll(s)
		if r.Type == Failure {
			g.o.WriteErr("Unknown code %v\n", s)
			continue
		}
		g.o.WriteAside("%s\n", r.Txt)
		if r.Type != Success {
			continue
		}
		if len(r.Url) != 0 {
			webhook(r.Url)
		}
		if r.Final {
			return
		}
	}
}

func (s *AnswerSet) Satisfies(g *PlayerState) bool {
	c := len(s.Requires)
	glog.Infof("looking for %v\n", c)
	for _, t := range s.Requires {
		glog.Infof("Checking %v\n", t)
		if slices.Contains(g.tags, t) {
			glog.Infof("got %v\n", t)
			c -= 1
		}
	}
	glog.Infof("missing %v\n", c)
	return c == 0
}

func (s *AnswerSet) Update(g *PlayerState) {
	g.tags = slices.DeleteFunc(g.tags, func(str Tag) bool {
		return slices.Contains(s.Consumes, str)
	})
	for _, t := range s.Provides {
		if !slices.Contains(g.tags, t) {
			g.tags = append(g.tags, t)
		}
	}
}

func (s *AnswerSet) Try(txt string, state *PlayerState) (r *Result) {
	r = &Result{Type: Failure, Txt: "", Final: false}
	transform := func(s string) string {
		return s
	}
	if !s.Strict {
		transform = func(s string) string {
			return strings.ToLower(strings.TrimSpace(s))
		}
	}

	for k, v := range s.Accept {
		if transform(k) == transform(txt) {
			if !s.Satisfies(state) {
				r.Type = PartialSuccess
				// TODO provide string
				r.Txt = "Almost...But you're missing something..."
				return
			}
			s.Update(state)
			r.Url = s.Webhook
			r.Type = Success
			r.Txt = v
			r.Final = s.Final
			return
		}
	}

	for k, v := range s.Partial {
		if transform(k) == transform(txt) {
			r.Type = PartialSuccess
			r.Txt = v
			return
		}
	}

	return r
}

func parse(model string) (s *gamemodel) {
	glog.Infof("Parsing...")
	s = &gamemodel{}

	err := yaml.Unmarshal([]byte(model), s)
	if err != nil {
		glog.Fatalf("Error unmarshalling model \"%v\": %v", model, err)
	}

	for i, state := range s.States {
		glog.Infof("State %v:\n", i)
		glog.Infof("\tAccept:\n")
		for k, v := range state.Accept {
			glog.Infof("\t\t%v -> %v\n", k, v)
		}
		glog.Infof("\tPartial:\n")
		for k, v := range state.Partial {
			glog.Infof("\t\t%v -> %v\n", k, v)
		}
		glog.Infof("\tRequires: %v\n", state.Requires)
		glog.Infof("\tProvides: %v\n", state.Provides)
		glog.Infof("\tConsumes: %v\n", state.Consumes)
		glog.Infof("\tStrict: %v\n", state.Strict)
		glog.Infof("\tFinal: %v\n", state.Final)
	}

	glog.Infof("Inital tags\n")
	for _, tag := range s.Tags {
		glog.Infof("\t%v\n", tag)
	}
	return
}
