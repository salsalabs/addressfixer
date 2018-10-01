package addressfixer

import (
	"log"
	"sync"
)

//After record changes to a supporter record.
func (e *Env) after(id int) {
	//log.Printf("After_%02d: start", id)
	for s := range e.After {
		e.DB.After(s)
	}
	log.Printf("After_%02d: end", id)
}

//Afterers starts tasks that record errors.
func (e *Env) Afterers(w *sync.WaitGroup, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, i int) {
			w.Add(1)
			e.after(i)
			w.Done()
		})(e, i)
	}
	log.Printf("Afterers: started %d Afters\n", c)
}
