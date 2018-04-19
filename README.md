# Awesome generator

![Awesome Generator logo](logo.png)

[![Go report](https://goreportcard.com/badge/github.com/orsinium/awesome-generator)](https://goreportcard.com/report/github.com/orsinium/awesome-generator) [![Code size](https://img.shields.io/github/languages/code-size/orsinium/awesome-generator.svg)](https://github.com/orsinium/awesome-generator) [![License](https://img.shields.io/github/license/orsinium/awesome-generator.svg)](LICENSE)

Generate awesome list for any language over [Github search API](https://developer.github.com/v3/search/#search-repositories).

Generated awesome lists: [generated-awesomeness](https://github.com/orsinium/generated-awesomeness).

## Usage

Generate awesome list for language:

```bash
go run main.go -l python > python.md
```

Generate awesome list for topic:

```bash
go run main.go -t docker > docker.md
```

## Advanced usage

Save projects to JSON:

```bash
go run main.go -l python --json > python.json
```

Generate awesome list from JSON:

```bash
cat python.json | go run main.go > python.md
```

## Command line arguments

* `-l` -- language. `go run main.go -l python`
* `-t` -- topic. `go run main.go -t docker`
* `--json` -- dump projects to json. `go run main.go -l python --json`
* `--pages` -- count of pages (default 10). `go run main.go -l python --pages 5`
* `--min` -- minimum projects into one section (default 2). `go run main.go -l python --min 3`
