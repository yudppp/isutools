package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/yudppp/isutools/tracereporter"
	goji "goji.io"
	"goji.io/pat"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var dbx *sqlx.DB

func main() {
	var err error
	dsn := getDSN()
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("mysql"))
	dbx, err = sqlxtrace.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
	}
	defer dbx.Close()

	tracereporter.Start(time.Second*15, tracereporter.NewHTTPRequestReport())
	mux := goji.NewMux()
	mux.Use(accessLoggeerSutpid())
	mux.HandleFunc(pat.Get("/hello/:name"), hello)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	user := map[string]interface{}{}
	dbx.GetContext(r.Context(), user, "SELECT * FROM `users`")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func accessLoggeerSutpid() func(next http.Handler) http.Handler {
	var urlReg = regexp.MustCompile(`([0-9]+)`)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			opts := []ddtrace.StartSpanOption{
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.ServiceName("router"),
			}
			span, ctx := tracer.StartSpanFromContext(r.Context(), "http.request", opts...)
			defer span.Finish()
			next.ServeHTTP(w, r.WithContext(ctx))
			urlName := urlReg.ReplaceAllLiteralString(r.URL.RequestURI(), ":id")
			resourceName := r.Method + " " + urlName
			span.SetTag(ext.ResourceName, resourceName)
		})
	}
}

func getDSN() string {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("failed to read DB port number from an environment variable MYSQL_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "isucari"
	}
	dbname := os.Getenv("MYSQL_DBNAME")
	if dbname == "" {
		dbname = "isucari"
	}
	password := os.Getenv("MYSQL_PASS")
	if password == "" {
		password = "isucari"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)
}
