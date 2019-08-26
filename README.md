# telnet
Telnet implementation in Go

## Installation

Run the following command from you terminal:

```bash
go get github.com/koind/telnet
```

## Usage

Usage example.

```
./telnet --address "site.com:80" --timeout 10000 -r 400
```

## View command line options

Run the following command from you terminal:
```
./telnet --help
```

Help information

```
Usage of ./telnet:
  -a, --address string    resource address for the connection
  -r, --read_timeout int  read timeout to connect
  -t, --timeout int       timeout to connect
```