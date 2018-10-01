package addressfixer

import (
	"github.com/salsalabs/godig"
)

//Supporter defines the parts of the supporter record that this app uses.
type Supporter struct {
	Key     string `json:"supporter_KEY"`
	Email   string
	City    string
	State   string
	Zip     string
	Country string
	TLD     string
}

//Loggable is a supporter and a message
type Loggable struct {
	S Supporter
	M string
}

//LoggableErr is a supporter and an error.
type LoggableErr struct {
	S Supporter
	E error
}

//Env is the runtime Environment.
type Env struct {
	Table  godig.Table
	Read   chan int32
	Fix    chan Supporter
	Save   chan Supporter
	Before chan Supporter
	After  chan Supporter
	Log    chan Loggable
	LogErr chan LoggableErr
	DB     *DBS
}

//DBLogin contains the information used to log into the database that the app
//uses.  The app supports any SQL databse.
//
//If you're using SQLite, then the type is "sqlite" and the filename is the
// path to the database.
//
//If you're using MySQL/MariaDB, then fill in the other fields.
type DBLogin struct {
	Type     string
	Database string
	User     string
	Password string
	Filename string
}
