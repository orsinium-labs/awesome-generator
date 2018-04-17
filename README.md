# Awesome generator

Generate awesome list for any language over [Github search API](https://developer.github.com/v3/search/#search-repositories).

Generated awesome lists: [generated-awesomeness](https://github.com/orsinium/generated-awesomeness).

## Usage

Let's generate awesome for Python!

Download and save data from Github API:

```bash
go run main.go -l python > python.json
```

Generate awesome list:

```bash
cat python.json | go run main.go > python.md
```

Keys:

* `-l` -- language. `go run main.go -l python`
* `-t` -- topic. `go run main.go -t monitoring`
* `--pages` -- count of pages (default 10). `go run main.go -l python --pages 5`
