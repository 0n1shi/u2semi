# U2semi

A web server running as a honeypot handling any requests.

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

## Configuration

### Custom header

You can add some items on response header like below.

```yaml
web:
   headers:
      - key: Server
        value: Apache/2.4.2 (Unix) PHP/4.2.2
```

### Custom content (json)

```yaml
web:
   contents:
      /greet:
         body: 'Hello world'
      /ping:
         body: '{"message":"pong"}'
```

### Custom content (directory)

This setting is given priority over json above.

```yaml
web:
   content_directory: ./content/ # directory path which has files to return as response
   directory_listing_template: ./template/directory_listing.html # html template for directory listing
```

### Request logging

Every request is saved to the repositories configured under `repo`.
If both are set, requests are saved to both. If neither is set, requests are only printed to stdout.

```yaml
repo:
   dsn: host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable # PostgreSQL
   file: ./log/requests.jsonl # JSON Lines file (one request per line, appended)
```

Each line of the file looks like below. `body_base64` is added only when the body is not valid UTF-8.

```json
{"received_at":"2026-10-06T18:53:32.786596+09:00","method":"POST","url":"/login","proto":"HTTP/1.1","headers":{"Content-Type":"application/x-www-form-urlencoded","User-Agent":"curl/8.7.1"},"body":"user=admin&pass=1234","ip_from":"192.0.2.1","ip_to":"198.51.100.1"}
```

Request bodies larger than 10 MiB are truncated.
