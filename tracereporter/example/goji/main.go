package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/yudppp/isutools/tracereporter"
	goji "goji.io"
	"goji.io/pat"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

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

func main() {
	tracereporter.Start(time.Second*15, tracereporter.NewSimpleReport(tracereporter.WithServiceName("router")))
	mux := goji.NewMux()
	mux.Use(accessLoggeerSutpid())
	mux.HandleFunc(pat.Get("/hello/:name"), hello)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}
