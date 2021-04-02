# icecream-go

🍦 *Never use `fmt.Println()` for debugging again! A port of the Python [IceCream](https://github.com/gruns/icecream) library*

---

### Features

* Close feature parity to the Python implementation
* Syntax highlighting
* Customisable outputs

### Inspecting variables

```go
// import "github.com/codemicro/icecream-go/ic"

foo := func (i int) int { return i + 333 }
ic.IC(foo(123)) // -> ic| foo(123): 456

bar := map[string]map[int]string{"key": {1: "one"}}
ic.IC(bar["key"][1]) // -> ic| bar["key"][1]: "one"

baz := struct{ Name string }{Name: "codemicro"}
ic.IC(baz.Name) // -> ic| baz.Name: "codemicro"
```

### Inspecting flows

```go
package main

import "github.com/codemicro/icecream-go/ic"

func foo() {
    ic.IC()
    // do things
    if condition {
        ic.IC() 
        // other thing 
    } else {
        ic.IC()
        // another thing
    }
}

func main() {
    foo()
    // -> ic| main.go:7 in github.com/codemicro/something/main.foo
    //    ic| main.go:9 in github.com/codemicro/something/main.foo
}
``` 

### `ic.Format`

```go
theAnswer := 42
x := ic.Format(theAnswer)
fmt.Println(x) // -> ic| theAnswer: 42
```
