# Http honeypot

A web server running as a honeypot handling any requests.

## Tech

- Go
- Gorm

## Usage

```bash
$ http-hoenypot
NAME:
   http honeypot - http server working as honeypot

USAGE:
   http honeypot [global options] command [command options] [arguments...]

COMMANDS:
   server, s   start honeyport http server
   version, v  show version
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
