package addressfixer

import (
	"log"
	"strings"
	"sync"
)

const (
	//Header goes at the start of the CSV output.
	Header string = "supporter_KEY,City,State,Zip,Country"
)

//Save stores the array of supporters in the database.
func (e *Env) save(id int) {
	//log.Printf("Save_%02d: start", id)
	for s := range e.Save {
		p := []string{
			"object=supporter",
			"key=" + s.Key,
			"City=" + s.City,
			"State=" + s.State,
			"Zip=" + s.Zip,
			"Country=" + s.Country}
		x := strings.Join(p, "&")
		x = strings.Replace(x, " ", "%20", -1)
		_, err := e.Table.Save(x, s.Key)
		if err != nil {
			panic(err)
		}
		e.After <- s
	}
	log.Printf("Save_%02d: end", id)
}

//Savers start the tasks that save to Salsa.
func (e *Env) Savers(w *sync.WaitGroup, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, id int) {
			w.Add(1)
			e.save(i)
			w.Done()
		})(e, i)
	}
	log.Printf("Savers: started %d Saves\n", c)
}
