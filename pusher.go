package addressfixer

import (
	"log"
	"strconv"
	"sync"
)

//Push computes the number of records based on the criteria, then
//pushes all of the offsets onto the read channel.  The channel is closed
//after that so that the readers will shut down neatly.
func (e *Env) Push(w *sync.WaitGroup, crit string) {
	w.Add(1)
	x, err := e.Table.Count(crit)
	if err != nil {
		log.Fatalf("Push: count error %v\n", err)
	}
	count, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		log.Fatalf("Push: count conversion error %v\n", err)
	}
	log.Printf("Push: processing %d supporters.\n", count)
	count = count + 499
	for i := int32(0); i < int32(count); i += 500 {
		e.Read <- i
	}
	close(e.Read)
	w.Done()
}
