# Streamlog

A real-time log streaming web application that reads logs from stdin and broadcasts them to multiple web clients via SSE, with live filtering capabilities.

## Requirements

- Go 1.24 or later
- Node.js 20 or later
- pnpm 8 or later

## Building

1. Build the Angular frontend:
```bash
go tool task build
```

## Testing

### Backend Tests
```bash
go tool ginkgo ./...
```

### Integration Tests
```bash
go tool ginkgo ./test/integration/...
```

## Usage

1. Start the application:
```bash
# Development mode
go tool task dev

# Production mode
./streamlog_go
```

2. Pipe logs to the application:
```bash
# Example: tail system logs
tail -f /var/log/syslog | ./streamlog_go

# Example: watch docker logs
docker logs -f container_name | ./streamlog_go
```

3. Open your browser and navigate to `http://localhost:<port>` (the port will be displayed in the console output)

### Features

- Real-time log streaming via Server-Sent Events (SSE)
- Live filtering of logs
- Multiple client support
- Automatic reconnection on connection loss

### Command Line Options

- `--port`: Specify the port to listen on (default: random available port)
- `--db`: Path to SQLite database file (default: in-memory database) 