package worker

import (
	"crypto/sha256"
	"fmt"

	"github.com/sonyabrk/image-palette-service/internal/cache"
	"github.com/sonyabrk/image-palette-service/internal/processor"
)

type Job struct {
	data     []byte
	resultCh chan<- jobResult
}

type jobResult struct {
	result *processor.Result
	err    error
}

type Pool struct {
	jobs  chan Job
	proc  *processor.Processor
	cache *cache.Cache
	done  chan struct{}
}

func NewPool(n int, proc *processor.Processor, c *cache.Cache) *Pool {
	p := &Pool{
		jobs:  make(chan Job, 100),
		proc:  proc,
		cache: c,
		done:  make(chan struct{}),
	}

	for i := 0; i < n; i++ {
		go p.runWorker()
	}

	return p
}

func (p *Pool) runWorker() {
	for job := range p.jobs {
		result, err := p.process(job.data)

		job.resultCh <- jobResult{result: result, err: err}
	}
}

func (p *Pool) process(data []byte) (*processor.Result, error) {
	key := hashImage(data)
	if cached, ok := p.cache.Get(key); ok {
		return cached.(*processor.Result), nil
	}

	result, err := p.proc.Analyze(data)
	if err != nil {
		return nil, err
	}

	p.cache.Set(key, result)
	return result, nil
}

func (p *Pool) Submit(data []byte) (*processor.Result, error) {
	resultCh := make(chan jobResult, 1)
	p.jobs <- Job{
		data:     data,
		resultCh: resultCh,
	}

	res := <-resultCh
	return res.result, res.err
}

func (p *Pool) Stop() {
	close(p.jobs)
}

func hashImage(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
