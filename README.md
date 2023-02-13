# U2semi

A web server running as a honeypot handling any requests.

## Tech

- Go
- Gorm

## Usage

```bash
NAME:
   U2semi - A honeypot working as a HTTP server

USAGE:
   U2semi [global options] command [command options] [arguments...]

COMMANDS:
   server, s   Start HTTP server
   version, v  Show version
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### Custom header

You can add items on the header like below.

```yaml
web:
  port: 80
  headers:
    - key: Server
      value: Apache/2.4.2 (Unix) PHP/4.2.2
```

### Custom content

```bash
# create a directory
mkdir db/

# keep empty dir as content
$ find ./content -type d -empty -exec touch {}/.gitkeep \;
```
