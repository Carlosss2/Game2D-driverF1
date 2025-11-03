package concurrency

import "sync"

// Job representa una tarea concurrente
type Job func() interface{}

func runJob(job Job, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	out <- job()
}


func forwardValues(in <-chan interface{}, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		out <- v
	}
}

func closeWhenDone(out chan interface{}, wg *sync.WaitGroup) {
	wg.Wait()
	close(out)
}



// FanOut ejecuta los jobs en goroutines y devuelve un canal con resultados
func FanOut(jobs []Job) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go runJob(job, out, &wg)
	}

	go closeWhenDone(out, &wg)
	return out
}

// FanIn combina varios canales en uno
func FanIn(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(chs))

	for _, ch := range chs {
		go forwardValues(ch, out, &wg)
	}

	go closeWhenDone(out, &wg)
	return out
}
