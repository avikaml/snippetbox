package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/avikaml/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// Application struct for holding the app-wide dependencies for the web app
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
}


func main() {

	// This is a command line flag. The default value is :4000
	addr := flag.String("addr", ":4000", "HTTP network adress")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// This reads in the command-line flag value and assigns it to  addr
	flag.Parse()

	// os.Stdout is the destination of the logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// New http.Server struct for the purpose of custom logging
	// By default, if Go’s HTTP server encounters an error it will log
	// it using the standard logger. For consistency it’d be better to
	// use our new errorLog logger instead.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	// the value returned from flag.String() is a pointer.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe() // We can call ListenAndServe because it satisfies the interface(?)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil{
		return nil, err
	}
	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}