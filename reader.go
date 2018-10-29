package addressfixer

import (
	"log"
	"sync"
)

//Read accepts an offset and reads 500 records (or maybe less)
//Records are sent to a downstream channel for processing.
func (e *Env) read(crit string, id int) {
	//log.Printf("Read_%02d: start", id)
	for offset := range e.Read {
		var a []Supporter
		count := 500
		includes := "supporter_KEY,City,State,Zip,Country"
		crit := crit + "&include=" + includes
		err := e.Table.Many(offset, count, crit, &a)
		// Errors are not fatal because deadline.
		if err != nil {
			log.Printf("Read_%02d: offset %7d %v\n", id, offset, err)
		} else {
			// Empty read returns [{	}].	Interesting, no?
			count = len(a)
			if count == 1 && len(a[0].Key) == 0 {
				count = 0
			}
			if count == 0 {
				log.Printf("Read_%02d: offset %7d, end of data\n", id, offset)
				return
			}
			for _, s := range a {
				e.Fix <- s
			}
		}
	}
	log.Printf("Read_%02d: end", id)
}

//Readers starts the tasks that read from Salsa.
func (e *Env) Readers(w *sync.WaitGroup, crit string, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, crit string, i int) {
			w.Add(1)
			defer w.Done()
			e.read(crit, i)
		})(e, crit, i)
	}
}
