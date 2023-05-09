package xros

import (
	"fmt"
	"io"
	"strings"

	"github.com/andrewz1/xros/proto"
)

type Client struct {
	cl io.Closer
	rd proto.Reader
	wr proto.Writer
}

func NewClient(rwc io.ReadWriteCloser) *Client {
	return &Client{
		cl: rwc,
		rd: proto.NewReader(rwc),
		wr: proto.NewWriter(rwc),
	}
}

func (c *Client) Close() error {
	return c.cl.Close()
}

func errWord(str string) error {
	return fmt.Errorf("unknown word: %s", str)
}

func (c *Client) readSentence() (*Sentence, error) {
	s := newSentence()
	for {
		str, err := c.rd.ReadWord()
		if err != nil {
			return nil, err
		} else if len(str) == 0 {
			break
		}
		switch str[0] {
		case '!':
			s.Name = str
		case '.':
			if !strings.HasPrefix(str, ".tag=") {
				return nil, errWord(str)
			}
			s.Tag = str[5:]
		case '=':
			t := strings.SplitN(str[1:], "=", 2)
			if len(t) == 1 {
				s.Data[t[0]] = ""
			} else {
				s.Data[t[0]] = t[1]
			}
		default:
			return nil, errWord(str)
		}
	}
	return s, nil
}

func (c *Client) Run(cmd ...string) ([]*Sentence, error) {
	err := c.wr.WriteWords(cmd...)
	if err != nil {
		return nil, err
	}
	ss := make([]*Sentence, 0, 10)
	var s *Sentence
	for {
		if s, err = c.readSentence(); err != nil {
			return nil, err
		} else if s.Name != SentenceRe {
			break
		}
		ss = append(ss, s)
		continue
	}
	if s.IsError() {
		return nil, s
	}
	if len(ss) == 0 {
		return nil, nil
	}
	return ss, nil
}

func (c *Client) Login(username, password string) error {
	_, err := c.Run("/login", "=name="+username, "=password="+password)
	return err
}
