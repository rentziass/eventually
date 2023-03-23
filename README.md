# Eventually

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rentziass/eventually)](https://pkg.go.dev/github.com/rentziass/eventually)

Eventually provides support for running a test block that should eventually succeed, while still
providing access to a `testing.TB` throughout the whole block. Eventually will keep trying re-running
the test block until either a timeout or max attempts are reached. While this was created to test
asynchronous systems I suppose it might help with flaky tests too :D. Here's an example:

```go
func TestAsync(t *testing.T) {
  asyncProcessSucceeded := false

  go func() {
    time.Sleep(100 * time.Millisecond)
    asyncProcessSucceeded = true
  }()

  eventually.Should(t, func(t testing.TB) {
    if !asyncProcessSucceeded {
      t.Fail()
    }
  })
}
```

Notice how within the function you pass to `eventually.Should` you have access to a `t testing.TB`, this
allows you to still use your helpers and assertions that rely on the good old `t`. Here's another example,
using [testify](https://github.com/stretchr/testify)'s `assert` and `require`:

```go
// code that sets up async consequences, e.g. writes some events on a queue

eventually.Should(t, func(t testing.TB) {
  events, err := readFromQueue()
  require.NoError(t, err)
  require.Len(t, events, 1)
  assert.Equal(t, "event", events[0])
})

```

Eventually has [`Should`](https://pkg.go.dev/github.com/rentziass/eventually#Should) and [`Must`](https://pkg.go.dev/github.com/rentziass/eventually#Must) functions, that 
correspond to [`Fail`](https://pkg.go.dev/testing#T.Fail) and [`FailNow`](https://pkg.go.dev/testing#T.FailNow) respectively in case of failure.

Behaviour can be customised with use of [`Options`](https://pkg.go.dev/github.com/rentziass/eventually#Option), for example:

```go
eventually.Should(t, func(t testing.TB) {
  // your test code here
},
  eventually.WithTimeout(10*time.Second),
  eventually.WithInterval(100*time.Millisecond),
  eventually.WithMaxAttempts(10),
)
```

And if you want to reuse your configuration you can do so by creating your very own `Eventually`. The example above would look something like:

```go
eventually := eventually.New(
  eventually.WithTimeout(10*time.Second),
  eventually.WithInterval(100*time.Millisecond),
  eventually.WithMaxAttempts(10),
)

eventually.Should(t, func(t testing.TB) {
  // test code
})

eventually.Must(t, func(t testing.TB) {
  // test code
})
```

## Why does this exist?

> TL;DR: I like `t` **a lot**

Other testing libraries have solutions for this. Testify for instance has its own [`Eventually`](https://pkg.go.dev/github.com/stretchr/testify@v1.8.2/assert#Eventually), but
the function it takes returns a `bool` and has no access to an "inner" `*testing.T` to be used for helpers and assertions.
Let's say for example that you have a helper function that reads a file and returns its content as a string, failing the test if 
it can't find the file (more convenient than handling all errors in the test itself). If the file you want to test is being
created asynchronously using that helper within Eventually will halt the whole test instead of trying executing again. In Go code:

```go
func TestAsyncFile(t *testing.T) {
  // setup

  assert.Eventually(t, func() bool {
    contents := readFile(t, "path") // <-- this halts the whole TestAsyncFile, not just this Eventually run
    return contents == "expected"
  })
}

func readFile(t *testing.T, path string) string {
  f, err := os.Open(path)
  require.NoError(t, err)

  // reading the file
}
```

Another available alternative is Gomega's [`Eventually`](https://pkg.go.dev/github.com/onsi/gomega#Eventually) (yes, this package has a very original name), which can be very convenient to use but requires buying into Gomega as a whole, which is quite the commitment (and I don't find a particularly idiomatic way of writing tests in Go but hey, opinions). This also still doesn't give access to a `t` with its own scope, you can do assertions within the `Eventually` block but if you have code that relies on `*testing.T` being around you cannot use it:

```go
gomega.Eventually(func(g gomega.Gomega) {
  contents := readFile(t, "path") // no t :(
  g.Expect(contents).To(gomega.Equal("expected"))
}).Should(gomega.Succeed())
```
