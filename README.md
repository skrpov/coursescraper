# UBCO Course scraper

Scrapes courses off of the ubco calendar site and presents them in graph form to 
make course dependencies more obvious.

## Requirments

- Go compiler, tested with `go version go1.20.5 windows/amd64`
- Graphviz, the program itself outputs a `.dot` file which Graphviz converts into an image.

## Building

```shell
$ go run .
$ dot -Tsvg ./save.dot > courses.svg
```