# Log
Go logging simplified.


```go
	var err error
	var l = log.V(0)

	// nil errors get ignored
	l.Err(err)

	// non nil errors are printed
	err = errors.New("we log non nil errors")
	l.Err(err)

	// or just write a string
	l.Str("hello")

	// and format it sometimes
	l.Strf("hello %s", "world")
```
