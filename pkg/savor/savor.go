package savor

import (
	"bytes"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// Processor handles the logic for what to do with matched requests for a Collector.
// It is created by library clients when setting up the routing matchers.
type Processor interface {
	Process(req *http.Request) error
}

// Collector is a container which handles any matched requests via the processor
type Collector struct {
	matcher   Matchable
	processor Processor
}

// NewCollector creates a new collector with the given matcher and processor
func NewCollector(m Matchable, p Processor) *Collector {
	return &Collector{m, p}
}

// consume returns true if it was matched and processed
func (c *Collector) Consume(req *http.Request) (bool, error) {
	if ok, err := c.matcher.matches(req); err != nil {
		return false, err
	} else if ok {
		c.processor.Process(req)
		return true, nil
	}
	return false, nil
}

// Matchable things implement a matches method for customization.
type Matchable interface {
	matches(*http.Request) (bool, error)
}

// GJsonMatcher is a matcher that handles matching request body via GJson
type GJsonMatcher struct {
	path  string
	equal any
}

func NewGJsonMatcher(path, value string) *GJsonMatcher {
	return &GJsonMatcher{path, value}
}

func (m *GJsonMatcher) matches(req *http.Request) (bool, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	if gjson.Get(string(body), m.path).String() == m.equal {
		return true, nil
	}
	return false, nil
}

// StringMatcher looks for the string anywhere in the body (as an example of another matcher)
type StringMatcher struct {
	value []byte
}

func NewStringMatcher(value string) *StringMatcher {
	return &StringMatcher{
		value: []byte(value),
	}
}

func (m *StringMatcher) matches(req *http.Request) (bool, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	res := bytes.Contains(body, m.value)
	return res, nil
}
