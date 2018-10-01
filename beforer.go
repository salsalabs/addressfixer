package addressfixer

import (
	"log"
	"sync"
)

//Before record changes to a supporter record.
func (e *Env) before(id int) {
	//log.Printf("Before_%02d: start", id)
	for s := range e.Before {
		e.DB.Before(s)
	}
	log.Printf("Before_%02d: end", id)
}

//Beforers starts tasks that record errors.
func (e *Env) Beforers(w *sync.WaitGroup, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, i int) {
			w.Add(1)
			e.before(i)
			w.Done()
		})(e, i)
	}
	log.Printf("Beforers: started %d Befores\n", c)
}
