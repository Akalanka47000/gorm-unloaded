# GORM Unloaded

A collection of self-contained Go examples covering the lesser-known features of [GORM](https://gorm.io) — the things that are mostly overlooked or buried in the corners of the docs.

## Prerequisites

- Go 1.25+
- Docker

## Getting started

Start the database containers (primary + replica):

```bash
make infra-up
```

Install dependencies:

```bash
go mod download
```

## Running an example

Each example is independent. Run any one by number:

```bash
make demo N=1
make demo N=5
make demo N=12
```

Each example creates and manages its own temporary database automatically — no setup required beyond the running containers.

## Project structure

```
.
├── db/                  # Shared database helpers — internal, ignore when reading examples
├── examples/
│   ├── 01-context/
│   ├── ...
│   └── 12-db-resolver/
├── docker-compose.yml   # PostgreSQL primary + replica (15432 / 15433)
└── go.mod
```
