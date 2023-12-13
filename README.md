# Dokka

## Overview

Dokka is a real-time monitoring and management tool for Docker containers.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Testing](#testing)
- [License](#license)

## Installation

Running Dokka

```bash
docker run --name dokka -d --volume=/var/run/docker.sock:/var/run/docker.sock -p 7070:7070 1704mori/dokka:latest
```
or with compose file
```bash
version: "3"

services:
  dokka:
    image: 1704mori/dokka:latest
    container_name: dokka
    ports:
      - 7070:7070
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

Dokka will be available at http://localhost:7070

## Configuration

Dokka can be configured using environment variables.

| Variable | Default | Description |
| --- | --- | --- |
| `DOKKA_PORT` | `7070` | Port to listen on |
| `DOKKA_EVENT_STREAM_INTERVAL` | `5` | Interval in seconds to send events to the client |

## Testing

```bash
go test ./...
```

## License

This project is licensed under the [MIT License](LICENSE).
