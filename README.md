# An embeddable Prolog

[ichiban/prolog](https://github.com/ichiban/prolog) embeddable scripting language for Go.

This package is intended to extend the `ichiban/prolog` with some helper functions and make it simple to calle Prolog from Golang.

### Usage

The package is fully go-getable, so, just type

  `go get github.com/rosbit/go-eprolog`

to install.

#### 1. Instantiate a Prolog interpreter

```go
package main

import (
  "github.com/rosbit/go-eprolog"
  "fmt"
)

func main() {
  ctx := pl.NewProlog()
  ...
}
```

#### 2. Load a Prolog script

Suppose there's a Prolog file named `music.pl` like this:

```prolog
listen(ergou, bach).
listen(ergou, beethoven).
listen(ergou, mozart).
listen(xiaohong, mj).
listen(xiaohong, dylan).
listen(xiaohong, bach).
listen(xiaohong, beethoven).
```

one can load the script like this:

```go
   if err := ctx.LoadFile("music.pl"); err != nil {
      // error processing
   }
```

#### 3. Prepare arguments and variables

```go
   // query Who listens to Music
   args := []interface{}{pl.PlVar("Who"), pl.PlVar("Music")}

   // query Who listens to "bach"
   args := []interface{}{pl.PlVar("Who"), "bach"}

   // query Which Music "ergou" listens to
   args := []interface{}{"ergou", pl.PlVar("Music")}

   // check whether "ergou" listens to "bach"
   args := []interface{}{"ergou", "bach"}
```

#### 4. Query the goal with arguments and variables

```go
   rs, ok, err := ctx.Query("listen", args...)
```

#### 5. Check the result

```go
   // error checking
   if err != nil {
      // error processing
      return
   }

   // proving checking with result `false`
   if !ok {
      // the result is false
      return
   }

   // proving checking with result `true`
   if rs == nil {
      // the result is true
      return
   }

   // result set processing
   for res := range rs {
      fmt.Printf("res: %#v\n", res)
   }
```

The full usage sample can be found [sample/main.go](https://github.com/rosbit/go-eprolog/blob/master/sample/main.go).

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.
__Convention:__ fork the repository and make changes on your fork in a feature branch.
