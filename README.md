GHA
======

GitHub Authentication.



Installation
---------

```sh
go get github.com/pocke/gha
```


Usage
------


### Basic usage


```go
import "github.com/pocke/gha"

func main() {
  key, err := gha.Auth("pocke", "PASSWORD", "Test App")
  if err != nil {
    panic(err)
  }
  fmt.Println(key)    # => Psersonal access token of GitHub
}
```

### For CLI Application

`main.go`

```go
import "github.com/pocke/gha"

func main() {
  key, err := gha.CLI(appName, "key.txt")
  if err != nil {
    panic(err)
  }
  fmt.Println(key)
}
```

Run

```sh
$ go run main.go
username: <INPUT YOUR USER NAME>
password for <YOUR USER NAME> (never stored): <INPUT YOUR PASSWORD>
<SHOW YOUR KEY>
```


`gha.CLI` saves your key to file.
If key is saved already, `gha.CLI` returns saved key.
