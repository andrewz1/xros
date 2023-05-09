package xros

import (
	"fmt"
	"strings"
)

const (
	SentenceDone = "!done"
	SentenceRe   = "!re"
	SentenceTrap = "!trap"
)

type Sentence struct {
	Name string
	Tag  string
	Data map[string]string
}

func newSentence() *Sentence {
	return &Sentence{
		Data: make(map[string]string),
	}
}

func (s *Sentence) String() string {
	if s == nil {
		return "sentence is nil"
	}
	sb := strings.Builder{}
	if len(s.Name) > 0 {
		sb.WriteString(s.Name)
	}
	if len(s.Tag) > 0 {
		fmt.Fprintf(&sb, " tag=%s", s.Tag)
	}
	for k, v := range s.Data {
		fmt.Fprintf(&sb, " %s=%s", k, v)
	}
	return sb.String()
}

func (s *Sentence) Error() string {
	return s.String()
}

func (s *Sentence) IsError() bool {
	return s == nil || s.Name != SentenceDone
}
