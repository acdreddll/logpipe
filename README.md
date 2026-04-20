# logpipe

A lightweight CLI tool for filtering and routing structured log streams in real time.

---

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git
cd logpipe
go build -o logpipe .
```

---

## Usage

Pipe any structured (JSON) log stream into `logpipe` and apply filters or route output to different destinations.

```bash
# Filter logs by level
my-app | logpipe --level error

# Filter by a specific field value
my-app | logpipe --filter service=auth

# Route logs to a file while still printing to stdout
my-app | logpipe --level warn --out warnings.log

# Combine filters
my-app | logpipe --level info --filter env=production
```

### Flags

| Flag | Description |
|------|-------------|
| `--level` | Minimum log level to display (`debug`, `info`, `warn`, `error`) |
| `--filter` | Filter by field value in `key=value` format (repeatable) |
| `--out` | Write matching logs to a file in addition to stdout |
| `--pretty` | Pretty-print JSON output |
| `--no-color` | Disable colored output when using `--pretty` |

---

## Example

```bash
$ echo '{"level":"error","msg":"connection refused","service":"db"}' | logpipe --level error --pretty
{
  "level": "error",
  "msg": "connection refused",
  "service": "db"
}
```

Multiple `--filter` flags can be combined (all must match):

```bash
my-app | logpipe --level warn --filter env=production --filter region=us-east-1
```

---

## License

MIT © 2024 yourusername
