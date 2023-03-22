# eventually

```go
// Must will keep trying until it succeeds or the test times out or reaches
// max attempts. If it fails, it will call t.Fatalf.
eventually.Must(t, func(t testing.TB) {
	t.Fatal("but keep trying")
}, eventually.WithMaxAttempts(10))

// Should will keep trying until it succeeds or the test times out or reaches
// max attempts. If it fails, it will call t.Errorf.
eventually.Should(t, func(t testing.TB) {
	t.Fatal("but keep trying")
}, eventually.WithTimeout(10*time.Second), eventually.WithInterval(1*time.Second))
```

## TODO
- [ ] Add documentation
