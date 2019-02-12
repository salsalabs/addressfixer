# Address Fixer for Salsa Classic

# Summary

This package contains an app that fixes cities, states and countries in Salsa Classic
supporter records.  The app accepts criteria for records to read.  Records are retrieved
from the database, cleaned up and stored back into the database.

This app uses a database to

-   Log non-fatal errors
-   Save the before-image of each modified record
-   Save the contents of each modified record

# How supporters are changed

-   If the country is a long name (like "Germany" or "France") then the Country field is changed to a [standard two-digit country code](http://www.nationsonline.org/oneworld/country_code_list.htm#).
-   The State field can be filled in using the Country and Zip fields.  (Zip in this
    case can also be postal codes outside the U.S.)
-   If the country is "CA" (Canada) and the state field is empty, then the state field is derived using the first
    digit of the Zip field (Canadian postal code). This is a workaround that solves the problem of Zippopatamus not beeing able to accurately retrieve
    all Canadian postal codes.  [Click here](https://en.wikipedia.org/wiki/Postal_codes_in_Canada) to the the article used to create the lookup table.

\*If the Country code is either empty or "US", then four-digit ZIP fields are changed to five digits by prepending a zero, but only if the states are in this list.

    * CT
    * MA
    * ME
    * NH
    * NJ
    * PR
    * RI
    * VT
    * VI

-   The City field can be filled in if the Country and Zip point to an unambigous city.

# Dependencies

This package is written in [Go](https://golang.org/).

This package depends on these apps:

-   [SQLLite v3](https://www.sqlite.org/index.html) Slick, powerful and server-less SQL.

This package depends on these Go packages.

-   [godig](https://github.com/salsalabs/godig) Go library for accessing Salsa Classic API.
-   [kingpin](https://github.com/alecthomas/kingpin) Go library to parse command-line arguments.
-   [go-sqlite3](https:github.com/mattn/go-sqlite3) SQLLite3 database driver used for testing
-   [go-mysql](https://github.com/go-sql-driver/mysql) MySQL database driver used in production

This package uses these web-bases services.

-   [RestCountries](https://restcountries.eu) to reduce long country names to ISO 2-digit country codes.
-   [Zippopotamus](http://zippopotam.us) to retrieve information about ZIP and postal codes.

# Build

Type this to create an executable file.

    go build -o addressfixer cmd/app/app.go
    mv addressfixer ~/go/bin

# Usage

Type this command to see the app-level help and usage.

    ./addressfixer --help

You will see something like this:

    usage: addressfixer --login=LOGIN [<flags>]

    Corrects cities, postal codes and countries in a Salsa database.

    Flags:
      --help         Show context-sensitive help (also try --help-long and --help-man).
      --login=LOGIN  YAML file with Salsa login credentials
      --dblogin      YAML file with database login information
      --criteria="Email IS NOT EMPTY&condition=Receive_Email>0&condition=State IS EMPTY&condition=Zip IS NOT EMPTY"
                     return supporters that match

The command-line arguments allow you to provide your Salsa Classic campaign manager credentials and modify the records that Salsa will retrieve.

## --login

A YAML file used to authenticate with your Salsa Classic HQ.  The field names
are self-explanatory.  Details can be found in the
[authenticate.sjs](https://help.salsalabs.com/hc/en-us/articles/115000341773#authenticatesjs)
documentation.  Here's a sample YAML file.

```yaml
host: https://salsa4.salsalabs.com
email: someone@somwhere.org
password: plaintext-of-your password
```

Needless to say, you'll need to make sure that no one else has access to the YAML
file.  Take the necessary precautions.  If you're not sure what to do, create
the file, run this app, then delete it.

## --dblogin

A YAML file used to specify the database and provide login credentials.  Here's a
sample of a dblogin file for MySQL:

    type: "mysql"
    login: "addressfixer:A_SAVE_PASSWORD@/addressfixer?charset=utf8"

Here's a sample of a dblogin file for SQLite:

    type: "sqlite3"
    filename: "./db/whatever.sqlite3"

## --criteria

The criteria tells addressfixer which records to retrieve from the database.  If
there are no criteria, then addressfixer will attempt to fix all supporters.  If
there is criteria, then it should be in the formation described in the
[condition](https://help.salsalabs.com/hc/en-us/articles/115000341773#condition)
documentation.

The criteria are provided to a call to
[getObjects.sjs](https://help.salsalabs.com/hc/en-us/articles/115000341773#getobjectsjs)
As a reminder, Salsa provides the first `&condition`.  If you need additional conditions,
then you'll need to separate them with `&condition`.

Here is an example of a working condition argument:

    'Email IS NOT EMPTY&condition=Receive_Email>0&condition=State IS EMPTY&condition=Zip IS NOT EMPTY'

This translates to parameters that retrieve a list of fields from active supporters that have empty State and postal code fields.

-   Email is not empty
-   Receive_Email > 0
-   State is empty
-   Zip is not empty

# Database

## Background

The application stores logging, before-image and after-image information in
a MySQL database named "addressfixer".  Application testing is done with a SQLite database.

Both databases run locally and can't be reached via the internets. Here's a summary of the tables used by addressfixer.

-   `log` contains the supporter record and a notation or an error
-   `afterimage` contains modified records
-   `preimage` contains the original contents of modified messages

## MySQL setup

If you are using MySQL/MariaDB, then you'll need to log in as root and use these commands:

    create database addressfixer;
    use addressfixer;
    grant all on addressfixer.* to addressfixer@localhost identified by "A-SAFE-PASSWORD";

Next, log in to MySQL/MariaDB as address fixer.  Use the contents of db/mysql_schema.sql to initialize the database.

The last step is to modify logins/db.yaml to store the login information.  Here's a sample:

    type: "mysql"
    database: "addressfixer"
    user: "addressfixer"
    password: "A_VERY_SAFE_PASSWORD"

The app will panic and die if the credentials are not correct.

## SQLite setup

create a new database and run the initialization SQL using a command line like this:

    sqlite3 ./db/whatever.sqlite3 < db/sqlite3_schema.sql

The last step is to modify logins/db.yaml to store the login information.  Here's a sample:

    type: "sqlite3"
    filename: "./db/whatever.sqlite3"

# Privacy

No personal data is stored by this application.  Records stored in the database only contain these fields.

-   supporter_KEY
-   City
-   State
-   Zip
-   Country

# License

Read the [LICENSE](./LICENSE).

# Questions?

Use the [Issues](https://github.com/salsalabs/addressfixer/issues)
at the top of the repository to ask questions or
report errors.  Don't waste your time by bothering Salsalabs Support.
