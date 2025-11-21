# letscrum

Open Source Lightweight Scrum/Agile Project Management System.

## Features

- **Organization Management**: Create and manage organizations.
- **Project Management**: Create and manage projects within organizations.
- **Sprint Management**: Plan and track sprints.
- **Work Item Management**: Manage Epics, Features, and Tasks.
- **User Management**: User authentication and role management.
- **API First**: Built with gRPC and REST (via gRPC Gateway).

## Tech Stack

- **Language**: Go 1.22+
- **RPC Framework**: gRPC
- **API Gateway**: gRPC-Gateway
- **Database ORM**: GORM (MySQL, PostgreSQL)
- **Configuration**: Viper
- **CLI**: Cobra
- **Authentication**: JWT

## Getting Started

### Prerequisites

- Go 1.22 or higher
- MySQL or PostgreSQL database

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/letscrum/letscrum.git
   cd letscrum
   ```

2. Build the project:

   ```bash
   make build
   ```

### Configuration

The application uses a configuration file located at `config/config.yaml` or `$HOME/.letscrum.yaml`.

Example configuration (`config/config.yaml`):

```yaml
server:
  http:
    addr: 0.0.0.0:8081
  grpc:
    addr: 0.0.0.0:9091

data:
  database:
    driver: mysql
    host: 127.0.0.1
    port: 3306
    database: letscrum
    user: root
    password: root
```

### Running the Server

To start the server (both gRPC and HTTP):

```bash
make run
# OR
./bin/letscrum server
```

## Development

### Generate API Code

If you modify the `.proto` files in `api/`, regenerate the Go code:

```bash
make api_gen
```

### Linting

Run the linter:

```bash
make lint
```

## Project Structure

- `api/`: Protocol Buffer definitions.
- `cmd/`: Application entry points.
- `config/`: Configuration files.
- `internal/`: Internal application logic (DAO, Service, Model, etc.).
- `pkg/`: Public libraries.
- `schema/`: Database schema.
- `Makefile`: Build and utility commands.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
