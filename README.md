# aiven

golang client and console cli for aiven

Status: project has very limited functionality and is under heavy development

### Installation

```bash
go get github.com/savaki/aiven/...
```

### Console Usage

sample input

```bash
aiven kafka topics --email XXX --password XXX --project XXX --service XXX
```

sample output

```json
[
  {
    "cleanup_policy": "compact",
    "partitions": 3,
    "replication": 3,
    "retention_hours": 72,
    "state": "ACTIVE",
    "topic_name": "sample"
  }
]
```

### Usage 

```bash
NAME:
   aiven - console interface to aiven

USAGE:
   aiven [global options] command [command options] [arguments...]

VERSION:
   SNAPSHOT

COMMANDS:
     kafka    kafka related commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```