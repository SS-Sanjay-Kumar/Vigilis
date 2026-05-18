# Vigilis

High-throughput log ingestion and analysis system built with Go and PostgreSQL.

Vigilis is designed to handle large-scale log streams with a focus on ingestion performance, concurrency, and efficient querying. The current implementation focuses on asynchronous processing, bulk persistence, and scalable database architecture for time-series log data.


## Features

- Asynchronous log ingestion using Go channels
- Producer-consumer worker architecture
- Batch processing with bulk inserts
- PostgreSQL connection pooling using `pgxpool`
- Time-based PostgreSQL table partitioning
- Optimized indexing strategy(BRIN) for large log datasets
- Structured logging for observability and debugging
- Clean modular backend architecture


## Architecture Overview

```text
Client Request(i.e services)
      │
      ▼
 Gin HTTP Server
      │
      ▼
 Buffered Go Channel
      │
      ▼
 Worker Goroutines
      │
      ▼
 In-Memory Batch Aggregation
      │
      ▼
 PostgreSQL Bulk Insert
```

## Tech Stack

- **Backend:** Go
- **Framework:** Gin
- **Database:** PostgreSQL
- **Driver:** pgx / pgxpool
- **Logging:** Structured Logging (`zap`)
- **Hot Reloading:** Air


## Current Implementation

### Concurrent Ingestion Pipeline

Incoming logs are accepted asynchronously using a producer-consumer model built with Go channels. The API returns immediately after queueing the log, allowing workers to process ingestion independently from request handling.

### Batch Processing

Logs are buffered in memory and inserted into PostgreSQL with the help of **PostgreSQL Copy Protocol** instead of single-row writes. This significantly reduces database overhead under high throughput.

### PostgreSQL Partitioning

Log tables are partitioned by timestamp using PostgreSQL declarative partitioning to improve scalability and query performance for large datasets.

### Optimized Query Performance

The project uses raw SQL with `pgx` and an BRIN indexing strategy optimized for large append-heavy log tables.


## Project Structure

```text
.
├── README.md
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   └── postgres.go
│   ├── handler
│   │   ├── health_handler.go
│   │   └── log_handler.go
│   ├── logger
│   │   └── logger.go
│   ├── models
│   │   └── log_entry.go
│   ├── repository
│   │   └── log_repository.go
│   └── worker
│       └── log_worker.go
└── tmp
    ├── build-errors.log
    └── main
```


## Getting Started

### Prerequisites

- Go
- PostgreSQL and pgx/pgxpool

### Run the Project

```bash
air
```


## Performance Goals

- Sustain high-throughput log ingestion
- Minimize write amplification
- Maintain fast query performance on large datasets
- Support scalable concurrent processing


## Roadmap

- [ ] AI-powered anomaly detection pipeline
- [ ] Redis-backed worker communication
- [ ] TensorFlow model integration
- [ ] Metrics and observability with Prometheus
- [ ] Containerized deployment setup
- [ ] Production-grade monitoring and documentation


## Design Decisions

### Why Asynchronous Processing?

Separating request handling from database writes improves ingestion throughput and reduces request latency during traffic spikes.

### Why Bulk Inserts?

Vigilis was built to handle 10,000 logs per second, leveraging PostgreSQL Copy Protocol we can inserts batches of logs into the database without compromising performance. Vigilis is build with high throughput in mind.

### Why Raw SQL?

To implement and use Respository Layer, adding abstraction, decoupling and centralization of data and data access.

### Why Partitioning?

Partitioning keeps large log datasets(potentially millions of rows) manageable and improves overall database performance.

### Why BRIN?

Vigilis is expected to ingest 10,000 logs per second, meaning millions of rows. Using a clutered or non clustered indexing here will result in very large index files(>30GB). This will fill up the server RAM and cause reduce overall performance of the database queries. 

Block Range Index has tenths of megabytes of index files to cover this huge volume of data. The trade-off made here is, having lossy index to combat server RAM overloading and overall performance degradation

---
## Status

Active development. Current focus is on building The anomaly model, a **Deep Autoencoder** in TensorFlow to detect abnormal logs i.e potential security risks.