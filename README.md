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

## Deployment

The steps below deploy U2semi as a systemd service on a Linux (amd64) server.

### 1. Build

Build a static binary. `-trimpath` removes local file paths from the binary and from log output.

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-X main.version=$(git describe --tags --always)" \
  -o dist/u2semi ./cmd/server
```

Alternatively, push a tag (e.g. `v0.4.0`) and GoReleaser builds a release on GitHub Actions.

### 2. Install

```bash
scp dist/u2semi <server>:/tmp/u2semi

# on the server
sudo install -m 755 -o root -g root /tmp/u2semi /opt/u2semi
sudo install -d -m 750 -o root -g root /var/log/http_honeypot # request log (may contain credentials sent by attackers)
```

Put the config file at `/etc/http_honeypot/config.yml` (see [Configuration](#configuration) and `config.yaml.example`), for example:

```yaml
repo:
  file: /var/log/http_honeypot/requests.jsonl
web:
  port: 80
  content_directory: /opt/contents/
  directory_listing_template: /etc/http_honeypot/directory_listing_template.html
```

### 3. systemd service

`/etc/systemd/system/http_honeypot.service`:

```ini
[Unit]
Description=HTTP honeypot

[Service]
ExecStart=/opt/u2semi server --config /etc/http_honeypot/config.yml
User=root
Group=root
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now http_honeypot
```

### 4. Log rotation

The log file grows without limit. Rotate it with logrotate. `copytruncate` works without restarting the service because the file is opened in append mode (requests arriving during the copy may be lost).

`/etc/logrotate.d/http_honeypot`:

```
/var/log/http_honeypot/requests.jsonl {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

### 5. Verify

```bash
systemctl is-active http_honeypot
curl -s http://127.0.0.1/ping
sudo tail -n 1 /var/log/http_honeypot/requests.jsonl
```

### Upgrade

Back up the current binary and config, replace the binary, then restart.

```bash
TS=$(date +%Y%m%d%H%M%S)
sudo cp -p /opt/u2semi /opt/u2semi.bak.$TS
sudo cp -p /etc/http_honeypot/config.yml /etc/http_honeypot/config.yml.bak.$TS
sudo install -m 755 -o root -g root /tmp/u2semi /opt/u2semi
sudo systemctl restart http_honeypot
```

To roll back, copy the `.bak.$TS` files back and restart the service.

Note: `repo.type` and `repo.mysql` used in old versions (<= 0.2.x) were removed in v0.3.0 and are ignored. Use `repo.dsn` (PostgreSQL) and/or `repo.file` instead.
