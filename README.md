# fileman — scheduled file cleanup service 🧹⏰

A tiny Go service that periodically deletes files older than a configured age from configured directories, driven by a cron schedule.

It’s ideal for pruning logs, temp files, and other time-bound artifacts. Runs as a single binary or in Docker and logs what it deletes.

---

## Features
- Cron-based scheduling
- Watch multiple directories, each with its own age threshold
- Simple JSON config file
- Docker enabled

---

## How it works (high level)
- On startup, the service loads a JSON config (CONFIG_PATH or `./config.json`).
- For each watched directory, it schedules a job using the configured cron expression.
- On each run, it lists entries in the directory and deletes files whose age (based on last modified time) is strictly greater than the threshold.
- It logs every deleted file and any errors; if nothing happens, it logs that too.

Notes:
- Age unit is days; fractional days are supported (e.g., 0.5 = 12 hours).
- Only files are deleted. Directories are ignored and not traversed.

---

## Configuration

```json
{
  "cron": "* * * * *",
  "watchedDirectories": [
    { "path": "/path/to/dir", "age": 7 }
  ]
}
```

Fields:
- cron: 5-field cron expression (minute precision). Example: `0 * * * *` = hourly at minute 0.
- watchedDirectories: array of objects with:
  - path: absolute path to the directory to prune
  - age: delete files older than this many days (float allowed)

---

## Quick start (local)
Prereqs: Go (per go.mod), or use Docker below.

1) Copy and edit config
```bash
cp config.json my-config.json
# edit my-config.json to add your directories
```

2) Run
```bash
CONFIG_PATH=my-config.json go run .
```

3) Build (optional)
```bash
go build -o fileman .
CONFIG_PATH=my-config.json ./fileman
```

4) Test
```bash
go test ./...
```

---

## Docker
Build a local image:
```bash
docker build -t fileman:local .
```

Run with bind mounts and non-root user:
```bash
docker run --rm \
  -e CONFIG_PATH=/app/config.json \
  -v $(pwd)/my-config.json:/app/config.json:ro \
  -v /path/on/host:/files \
  fileman:local
```

- PUID/PGID build args are supported to set the container user (defaults 1000/1000). If needed, build like:
```bash
docker build --build-arg PUID=$(id -u) --build-arg PGID=$(id -g) -t fileman:local .
```

---

## Docker Compose
A compose file is included. Edit the volume paths for your system and config.

Example:
```yaml
services:
  fileman:
    build:
      context: .
      args:
        - "PUID=${PUID:-1000}"
        - "PGID=${PGID:-1000}"
    image: murilopereira.dev/fileman:latest
    container_name: fileman
    restart: unless-stopped
    environment:
      - CONFIG_PATH=/app/config.json
    volumes:
      - /absolute/host/folder:/files
      - ./my-config.json:/app/config.json:ro
```
Start it:
```bash
docker compose up -d --build
```

---

## Safety and limitations
- One level only: does not recurse into subdirectories.
- Directories are never removed; only files can be deleted.
- Deletions are permanent. Review your config carefully and test on a sample directory first.
- File age uses last modified time (mtime).
- If a directory is unreadable or a file can’t be removed, the error is logged and processing continues.

---

## Environment
- CONFIG_PATH: optional; path to the JSON config (default `config.json` in working directory). In Docker Compose we use `/app/config.json`.

---

## Development
- Run tests: `go test ./...`
- Project layout: small, modular packages: clock, config, fs, handler
- Scheduler: github.com/go-co-op/gocron/v2

---

## Example config
```json
{
  "cron": "0 * * * *",
  "watchedDirectories": [
    { "path": "/files/logs", "age": 7 },
    { "path": "/files/tmp",  "age": 0.5 }
  ]
}
```

This runs hourly, deleting files older than 7 days in `/files/logs` and older than 12 hours in `/files/tmp`.

---

## License

This project is licensed under MIT License. See [LICENSE](LICENSE) for details.

---

Thanks for using **fileman**! If you need help or have ideas to improve it, feel free to raise issues or submit pull requests.