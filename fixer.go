package addressfixer

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

//Fix updates a supporter record using a couple of free sites.
func (e *Env) fix(id int) {
	noZip := regexp.MustCompile("(ZIP+|UNDEF+)")

	//log.Printf("Fix_%02d: start", id)
	for s := range e.Fix {
		modified := false

		before := Supporter{
			Key:     strings.TrimSpace(s.Key),
			City:    strings.TrimSpace(s.City),
			State:   strings.TrimSpace(s.State),
			Zip:     strings.TrimSpace(s.Zip),
			Country: strings.TrimSpace(s.Country),
		}
		// Move this after modified check if we're only
		//keeping before-images for modified records.
		e.Before <- before

		// Skip the usual suspects.
		z := strings.ToUpper(s.Zip)
		if noZip.MatchString(z) {
			m := Loggable{
				S: s,
				M: "'ZIP' in Zip",
			}
			e.Log <- m
		} else {
			// Get country code for long country name.
			// Do this before jumping into the postal code lookup.
			m, err := RestCountries(&s)
			if err != nil {
				m := LoggableErr{
					S: s,
					E: err,
				}
				e.LogErr <- m
			}
			modified = modified || m
			m, err = Zippo(&s)
			if err != nil {
				m := LoggableErr{
					S: s,
					E: err,
				}
				e.LogErr <- m
			}
			modified = modified || m
			msg := fmt.Sprintf("Fix_%02d: fixed? %v", id, modified)
			x := Loggable{
				S: s,
				M: msg,
			}
			e.Log <- x

			if modified {
				//e.Before <- before
				e.Save <- s
			}
		}
	}
	log.Printf("Fix_%02d: end", id)
}

//Fixers starts the Fix tasks.
func (e *Env) Fixers(w *sync.WaitGroup, c int) {
	for i := 0; i < c; i++ {
		go (func(e *Env, i int) {
			w.Add(1)
			e.fix(i)
		})(e, i)
	}
}
