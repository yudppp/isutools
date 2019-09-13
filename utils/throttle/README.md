# throttle

throtting function called
sync.Once like interface

```go
var profileThrottle = throttle.New(time.Second*5)

func SomeFunc() {
    profileThrottle.Do(func(){
        fmt.Println("run")
    })
}
```

