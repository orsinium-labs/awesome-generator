# Awesome generator

Generate awesome list for any language over [Github search API](https://developer.github.com/v3/search/#search-repositories).

Generated awesome lists: [generated-awesomeness](https://github.com/orsinium/generated-awesomeness).

## Usage

Let's generate awesome for Python!

Generate awesome list:

```bash
go run main.go -l python > python.md
```

Download and save data from Github API:

```bash
go run main.go -l python -j > python.json
```

Keys:

* `-l` -- language. `go run main.go -l python`
* `-t` -- topic. `go run main.go -t monitoring`
* `-j` -- dump json. `go run main.go -l python -j`
* `--pages` -- count of pages (default 10). `go run main.go -l python --pages 5`
