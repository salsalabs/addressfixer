package addressfixer

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	//Need this to drive the SQLite3 database.
	_ "github.com/mattn/go-sqlite3"
	//Need this to drive the MySQL database.
	_ "github.com/go-sql-driver/mysql"
)

//DBS is the struct that hold the app database stuff.
type DBS struct {
	D            *sql.DB
	LogStmt      *sql.Stmt
	PostStmt     *sql.Stmt
	PreStmt      *sql.Stmt
	LastPostStmt *sql.Stmt
}

//NewDBS returns a app database object.  The parameter is the name of a YAML
//file that contains the login parameters.
func NewDBS(cpath string) (*DBS, error) {
	raw, err := ioutil.ReadFile(cpath)
	if err != nil {
		return nil, err
	}
	var c DBLogin
	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DB credentials are %+v\n", c)
	dbArg := ""
	switch c.Type {
	case "sqlite3":
		dbArg = c.Filename
	case "mysql":
		t := "%v:%v@/%v?charset=utf8"
		dbArg = fmt.Sprintf(t, c.User, c.Password, c.Database)
	default:
		err := fmt.Errorf("'%v' not a valid database type", c.Type)
		return nil, err
	}
	db := DBS{}
	x, err := sql.Open(c.Type, dbArg)
	if err != nil {
		return &db, err
	}
	db.D = x

	s, err := db.D.Prepare("insert into log(id,city,state,zip,country,reason) values(?,?,?,?,?,?)")
	if err != nil {
		return &db, err
	}
	db.LogStmt = s

	s, err = db.D.Prepare("insert into preimage(id,city,state,zip,country) values(?,?,?,?,?)")
	if err != nil {
		return &db, err
	}
	db.PreStmt = s

	s, err = db.D.Prepare("insert into postimage(id,city,state,zip,country) values(?,?,?,?,?)")
	if err != nil {
		return &db, err
	}
	db.PostStmt = s

	s, err = db.D.Prepare("select at from postimage order by at desc limit 1;")
	if err != nil {
		return &db, err
	}
	db.LastPostStmt = s

	return &db, err
}

//After puts a supporter into after-image table.
func (db *DBS) After(s Supporter) error {
	_, err := db.PostStmt.Exec(s.Key, s.City, s.State, s.Zip, s.Country)
	if err != nil {
		return err
	}
	return nil
}

//Before saves a before image into the database.
func (db *DBS) Before(s Supporter) (err error) {
	_, err = db.PreStmt.Exec(s.Key, s.City, s.State, s.Zip, s.Country)
	return err
}

//Log writes a supporer record and a notation to the database.
func (db *DBS) Log(s Loggable) {
	_, err := db.LogStmt.Exec(s.S.Key, s.S.City, s.S.State, s.S.Zip, s.S.Country, s.M)
	if err != nil {
		//log.Fatalf("%v\n", err)
		log.Printf("Log: error %v on %+v\n", err, s)
		//panic(err)
	}
}

//LogErr writes an error to the database if err is not empty.
func (db *DBS) LogErr(s LoggableErr) {
	x := Loggable{
		S: s.S,
		M: fmt.Sprintf("%v", s.E),
	}
	db.Log(x)
}

//LastPost returns the date of the most recent postimage record.
func (db *DBS) LastPost() (at string, err error) {
	row := db.LastPostStmt.QueryRow()
	err = row.Scan(&at)
	return at, err
}
