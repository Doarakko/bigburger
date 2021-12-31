# Development Guide

Add new toppingss.

## 1. Set

```
git clone https://github.com/Doarakko/bigburger
cd bigburger
```

## 2. Prepare new topping image

Width is `218` pixels and format is `PNG`.

More than 218 pixels will be cut.

Set to `public/<image file name>.png`.

## 3. Generate `statik/statik.go`

```bash
statik -src=public
```

## 4. Fix `main.go`

Use local `statik`.

```go
import (
	// omit

	_ "./statik"
	//_ "github.com/Doarakko/bigburger/statik"

	// omit
)
```

Change array size.

```go
var toppings [n]Topping
```

Add some description about new topping.

Image file name and `toppings[i].Name` is same.
`toppings[i].Count` is default value.

```go
func init() {
	// omit

	toppings[6].Name = "<image file name>"
	toppings[6].Count = 0
	toppings[6].Option = "<option initial>"

	// omit
}
```

## 5. Run

```go
go run main.go -<option initial> 3
```
