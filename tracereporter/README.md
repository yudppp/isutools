# tracereporter

datadogのAPMの結果みたいなものをSlackに通知する
initial系のものにhooksさせて開始するようにする

```go
func main() {
    // ...
    // use https://github.com/DataDog/dd-trace-go package
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("mysql"))
	dbx, err = sqlxtrace.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
	}
	defer dbx.Close()
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
    }
    // ...
}
    
func postInitialize(w http.ResponseWriter, r *http.Request) {
    tracereporter.Start(benchTime, tracereporter.NewSimpleReport(tracereporter.WithServiceName("mysql")))
	// ...
}
```