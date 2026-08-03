# Analytics 77

**Analytics 77** is a high-performance, real-time analytics ingestion and aggregation engine written in Go. It is designed to efficiently collect, process, and aggregate time-series event data with minimal latency, zero-allocation patterns, and high concurrency.

> **License Notice**: This project is licensed under a proprietary Commercial License. You may view and test the code, but commercial use or distribution requires a paid license.

## 🚀 Key Features

- **Dual-Transport Ingestion**: Supports concurrent data ingestion via high-performance HTTP (Fiber v3) and custom TCP transports.
- **Time-Series Aggregation**: Optimized domain models (`AnalyticsCube`, `AnalyticsRegistry`) for rolling time-based aggregations (hourly, daily, monthly).
- **Timezone & DST Aware**: Built-in handling of UTC offsets and Daylight Saving Time (DST) transitions per site registry, ensuring accurate civil-calendar math at boundaries.
- **High-Performance Primitives**: Leverages atomic operations, arena allocators, and fine-grained `RWMutex` locking to maximize throughput and minimize garbage collection pressure.
- **Modular Architecture**: Clean separation of concerns across entrypoints, business logic, infrastructure, and shared utilities.

## 🏗️ Architecture Overview

- **`cmd/`**: Application entrypoints (e.g., `cmd-analytics`, `cmd-debug`).
- **`app-analytics/`**: Core application orchestration, dependency injection, and routing setup.
- **`domain/`**: Business logic and data structures, including analytics cubes, registries, and time-based aggregation rules.
- **`infra/`**: Infrastructure concerns, including the `DataCenter` (in-memory data management), TCP/HTTP transport layers, and geolocation services.
- **`services/`**: High-level services that bridge the infrastructure and domain layers.
- **`helpers/` & `shared/`**: Cross-cutting utilities, timestamp extraction, and shared constants.

## 🛠️ Getting Started

### Prerequisites

- Go 1.26 or higher


