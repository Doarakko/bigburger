# bigburger

[![Go Reference](https://pkg.go.dev/badge/github.com/Doarakko/bigburger.svg)](https://pkg.go.dev/github.com/Doarakko/bigburger)
[![Go Report Card](https://goreportcard.com/badge/github.com/Doarakko/bigburger)](https://goreportcard.com/report/github.com/Doarakko/bigburger)

Biiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiig burger!

![sample](./sample.gif)

## Install

go

```sh
go install github.com/Doarakko/bigburger@latest
```

Homebrew

```sh
brew install Doarakko/tap/bigburger
```

## Usage

```bash
Usage of bigburger:
  -C int
    	how many cat
  -b int
    	how many bun
  -c int
    	how many cheese
  -i float
    	rate of intervals (default 1)
  -l int
    	how many lettuce
  -n int
    	number of big burger (default 1)
  -o string
    	output image file
  -p int
    	how many putty (default 1)
  -s	buns with sesame
  -t int
    	how many tomato
```

## Example

```bash
bigburger -n 3 -p 6 -c 5 -l 3 -t 3
```

![example](./sample2.gif)

## Release

```sh
# 1. Create and push tag
git tag v0.1.0
git push origin v0.1.0
# 2. Update homebrew-tap
./scripts/update-homebrew.sh 0.1.0
```

## Credit

- [longcat](https://github.com/mattn/longcat)

## Development

Please check the [development guide](./DEVELOPMENT.md).

## Contribution

Let's add your favorite toppings!

## License

MIT

## Author

Doarakko
