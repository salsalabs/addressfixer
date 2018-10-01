package addressfixer

import (
	"log"
	"sync"
)

//Log record changes to a supporter record.
func (e *Env) log(id int) {
	//log.Printf("Log_%02d: start", id)
	for s := range e.Log {
		e.DB.Log(s)
	}
	log.Printf("Log_%02d: end", id)
}

//Log record changes to a supporter record.
func (e *Env) logErr(id int) {
	//log.Printf("LogErr_%02d: start", id)
	for s := range e.LogErr {
		e.DB.LogErr(s)
	}
	log.Printf("LogErr_%02d: end", id)
}

//Loggers starts tasks that record errors.
func (e *Env) Loggers(w *sync.WaitGroup, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, i int) {
			w.Add(1)
			e.log(i)
			w.Done()
		})(e, i)
	}
	log.Printf("Loggers: started %d Logs\n", c)

	for i := 0; i < c; i++ {
		go (func(e *Env, i int) {
			w.Add(1)
			e.logErr(i)
			w.Done()
		})(e, i)
	}
	log.Printf("Loggers: started %d LogErrs\n", c)
}
