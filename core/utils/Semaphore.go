package utils

type Empty struct {
}

// sem = make(semaphore, N)
type Semaphore chan Empty

/* mutexes */

func (s Semaphore) Lock() {
	s.acquire(1)
}

func (s Semaphore) Unlock() {
	s.release(1)
}

/* signal-wait */

func (s Semaphore) Signal() {
	s.release(1)
}

func (s Semaphore) Wait(n int) {
	s.acquire(n)
}

// acquire n resources
func (s Semaphore) acquire(n int) {
	e := Empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

// release n resources
func (s Semaphore) release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}
