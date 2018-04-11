

## Usage

Let's generate awesome for Python!

Download and save data from Github API:

```bash
go run main.go -l python > python.json
```

Generate awesome list:

```bash
cat python.json | go run main.go > awesomes/python.md
```
