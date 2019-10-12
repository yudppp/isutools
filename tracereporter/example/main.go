package main


import (
	"time"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/yudppp/isutools/tracereporter"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

func main() {
	tracereporter.Start(time.Second * 3, "", func(span mocktracer.Span) string {
		tags := span.Tags()
		resoureName, ok := tags["resource.name"]
		if !ok {
			fmt.Println(span.OperationName(), tags)
			return ""
		}
		return fmt.Sprint(resoureName)
	})

	dsn := getDSN()
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("mysql"))
	dbx, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
	}
	defer dbx.Close()

	user := map[string]interface{}{}
	dbx.Get(user, "SELECT * FROM `users`")

	item := map[string]interface{}{}
	dbx.Get(item, "SELECT * FROM `items`")

	time.Sleep(time.Second * 5)
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