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
