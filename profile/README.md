# profile

pprofした結果をSlackに通知する仕組み

initial系のものにhooksさせて開始するようにする

```go
func postInitialize(w http.ResponseWriter, r *http.Request) {
    profile.StartAll(time.Minute)
	// ...
}
```