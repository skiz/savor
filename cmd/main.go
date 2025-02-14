package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/skiz/savor/pkg/savor"
)

// NthProcessor triggers a printline on the 10th message received.
type NthProcessor struct {
	nth int32
	cur atomic.Int32
	mu  sync.RWMutex
}

// process implements savor.Processor.
func (p *NthProcessor) Process(req *http.Request) error {
	if p.cur.Add(1) == p.nth {
		p.cur.Add(-p.nth)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		fmt.Printf("Got %v requests. Body: %s", p.nth, body)
	}
	return nil
}

// Here is a timed count processor which prints out the total and resets a counter after a set amount of time.
type TimedCountProcessor struct {
	ticker time.Ticker
	done   chan bool
	cur    atomic.Int32
}

// trigger is called on every duration of the ticker
func (p *TimedCountProcessor) trigger() {
	n := p.cur.Load()
	p.cur.Add(-n)
	fmt.Println("TEST count:", n)
}

// stop isn't really used in this example
func (p *TimedCountProcessor) stop() {
	p.ticker.Stop()
	p.done <- true
}

// start the ticker
func (p *TimedCountProcessor) start() {
	go func() {
		for {
			select {
			case <-p.ticker.C:
				p.trigger()
			case <-p.done:
				return
			}
		}
	}()
}

// process implements savor.Processor. Just adds 1
func (p *TimedCountProcessor) Process(req *http.Request) error {
	p.cur.Add(1)
	return nil
}

// NewTimedCountProcessor configures and fires up the ticker.
func NewTimedCountProcessor(every time.Duration) *TimedCountProcessor {
	p := &TimedCountProcessor{
		ticker: *time.NewTicker(every),
		done:   make(chan bool),
	}
	p.start()
	return p
}

// Router is an example implementation of a HTTP router that just iterates
// a collection of matchers for the first successful match. Returns OK on valid requests.
// This could also be middleware, but we will just use a simple router for now.
type Router struct {
	collectors []*savor.Collector
}

func NewRouter(collectors ...*savor.Collector) *Router {
	return &Router{
		collectors: collectors,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, c := range r.collectors {
		if handled, err := c.Consume(req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 INTERNAL SERVER ERROR"))
			fmt.Println("unable to consume request:", req, err)
			return
		} else if handled {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 OK"))
			return // stop if we matched
		}
	}

	// If nothing matched, we just 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func main() {

	// We want to match the 10th Dave we see and print out a line.
	daveMatcher := savor.NewGJsonMatcher("name.first", "Dave")
	tenthProcessor := &NthProcessor{nth: 10}
	tenthDaveCollector := savor.NewCollector(daveMatcher, tenthProcessor)

	// We also want to match 'TEST' anywhere in the req body and print the count every 60 seconds
	testMatcher := savor.NewStringMatcher("TEST")
	timedCountProcessor := NewTimedCountProcessor(60 * time.Second)
	timedTestCollector := savor.NewCollector(testMatcher, timedCountProcessor)

	router := NewRouter(tenthDaveCollector, timedTestCollector)

	mux := http.NewServeMux()
	mux.Handle("/", router)
	http.ListenAndServe(":8080", mux)
}
