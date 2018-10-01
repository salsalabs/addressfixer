package main

import (
	"log"
	"os"
	"sync"

	"github.com/salsalabs/addressfixer"
	"github.com/salsalabs/godig"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const numFixers = 5
const numBeforers = 1
const numSavers = 2
const numAfters = 1
const numLoggers = 1
const numReaders = 5
const defCrit = `Email IS NOT EMPTY&condition=Receive_Email>0&condition=State IS EMPTY&condition=Zip IS NOT EMPTY`

//const dbType = "sqlite3"
//const dbArg = "db/addressfixer.sqlite3"
const dbType = "mysql"
const dbArg = "addressfixer:VAouGnSoEheBdmPYM9eTFjGKT7VXA9a@/addressfixer?charset=utf8"

//Fatal handles a fatal error.
func Fatal(err error) {
	log.Fatalf("%v\n", err)
}

func main() {
	var (
		app     = kingpin.New("addressfixer", "Corrects cities, postal codes and countries in a Salsa database.")
		login   = app.Flag("login", "YAML file with Salsa campaign manager credentials").Required().String()
		dbLogin = app.Flag("dblogin", "YAML file with database login credentials").Required().String()
		crit    = app.Flag("criteria", "return supporters that match").Default(defCrit).String()
	)
	app.Parse(os.Args[1:])
	api, err := (godig.YAMLAuth(*login))
	if err != nil {
		log.Fatalf("Main: authentication error %v\n", err)
	}

	table := api.NewTable("supporter")
	read := make(chan int32)
	fix := make(chan addressfixer.Supporter)
	before := make(chan addressfixer.Supporter)
	save := make(chan addressfixer.Supporter)
	after := make(chan addressfixer.Supporter)
	logr := make(chan addressfixer.Loggable)
	loge := make(chan addressfixer.LoggableErr)
	d, err := addressfixer.NewDBS(*dbLogin)
	if err != nil {
		panic(err)
	}

	e := addressfixer.Env{
		Table:  table,
		Read:   read,
		Fix:    fix,
		Save:   save,
		Before: before,
		After:  after,
		Log:    logr,
		LogErr: loge,
		DB:     d,
	}

	var w sync.WaitGroup
	e.Loggers(&w, numLoggers)
	e.Afterers(&w, numAfters)
	e.Beforers(&w, numBeforers)
	e.Savers(&w, numSavers)
	e.Fixers(&w, numFixers)
	e.Readers(&w, *crit, numReaders)
	e.Push(&w, *crit)
	w.Wait()
}
