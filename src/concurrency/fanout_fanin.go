package concurrency

import "sync"

// Job representa una tarea concurrente
type Job func() interface{}

// FanOut ejecuta los jobs en goroutines y devuelve un canal con resultados
func FanOut(jobs []Job) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			out <- j()
		}(job)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// FanIn combina varios canales en uno 
func FanIn(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, ch := range chs {
		go func(c <-chan interface{}) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
