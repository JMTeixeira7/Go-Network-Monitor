# HTTP Proxy in Go

A simple HTTP proxy written in Go.

This project uses a MySQL database running in Docker. The database is automatically initialized with an `init.sql` script through the Docker configuration.

## Requirements

- Go installed
- Docker installed
- Docker Compose installed

## Project setup

Clone the repository and enter the project folder:

```bash
git clone <repo-url>
cd <repo-folder>
```

### Start MySQL container:

Terminal1
```bash
cd <repo-folder>
docker compose --env-file .env -f internal/db/docker-compose.yml down -v
docker compose --env-file .env -f internal/db/docker-compose.yml up -d mysql
```
```bash
docker compose --env-file .env -f internal/db/docker-compose.yml ps
```

Terminal2
```bash
cd <repo-folder/internal/db>
docker compose --env-file ../../.env -f docker-compose.yml exec mysql mysql -uproxyuser -pproxypass proxydb
```

Terminal3
```bash
cd <repo-folder>
make run
```

### Set Up proxy on your Browser:

Got to settings and set a proxy route at:
```bash
host: 127.0.0.1
port: 4444
```