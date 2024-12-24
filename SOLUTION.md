# User Segmentation Challenge Solution

## Problem Statement

We need to build a system that:

1. Receives user segment data (user_id, segment pairs)
2. Maintains segments for users for two weeks
3. Can estimate the number of users in each segment
4. Must handle millions of users and hundreds of segments

## Solution Architecture

### 1. High-Level Design

```
[USS Service] -> [Kafka] -> [ES Service (Batch Processor)] -> [ClickHouse]
```

### 2. Components

#### 2.1 User Segmentation Service (USS)

- Receives segment data via API
- Produces messages to Kafka
- Stateless service for high scalability
- No direct database interaction

#### 2.2 Kafka Message Broker

- Buffers incoming segment data
- Ensures reliable message delivery
- Enables system scalability
- Decouples USS from ES

#### 2.3 Estimation Service (ES)

- Consumes messages from Kafka
- Implements batch processing
- Manages data persistence
- Handles estimation queries

#### 2.4 ClickHouse Database

- Stores user segment data
- Optimized for analytical queries
- Efficient at counting distinct users
- Handles time-based data expiry

### 3. Key Implementation Details

#### 3.1 Batch Processing Strategy

To handle 2 million connections efficiently, the batch processing strategy must be highly configurable. The following parameters can be adjusted to optimize performance:

- **Buffer Size**: Configurable to handle up to 2 million messages. Default is 1000 messages.
- **Flush Interval**: Configurable time interval for flushing the buffer. Default is 30 seconds.
- **Retry Attempts**: Configurable number of retry attempts for failed batches. Default is 3 attempts.

These parameters ensure that the system can scale and handle high volumes of data efficiently while maintaining reliability and performance.

#### 3.2 Data Flow

1. USS receives segment data
2. Data is published to Kafka
3. ES batch processes messages
4. Batches are stored in ClickHouse
5. Queries read from ClickHouse

#### 3.3 Clean Architecture Implementation

```
└── internal/
    ├── domain/        # Core business entities and interfaces
    ├── usecase/       # Application business logic
    ├── delivery/      # Interface adapters (e.g., API handlers)
    └── infrastructure/# Frameworks and drivers (e.g., database, messaging)
```

This structure ensures a clear separation of concerns, making the system more maintainable and testable.

### 4. Design Decisions

#### 4.1 Why Kafka?

Kafka is well-suited for processing big data due to its high throughput, fault tolerance, and scalability. It can handle large volumes of data by distributing the load across multiple brokers and partitions. This ensures that the system can ingest and process millions of user segment messages efficiently.

Key benefits of using Kafka for big data processing include:

- **High Throughput**: Capable of handling thousands of messages per second, making it ideal for big data applications.
- **Scalability**: Easily scales horizontally by adding more brokers and partitions.
- **Fault Tolerance**: Replicates data across multiple brokers to ensure reliability and data durability.
- **Real-Time Processing**: Supports real-time data streaming, enabling timely insights and actions.

By leveraging Kafka, our system can efficiently manage and process the large volumes of user segment data required for this challenge.

- Message buffering for scalability
- Reliable message delivery
- Supports multiple consumers
- Built-in partitioning

#### 4.2 Why ClickHouse?

- Optimized for analytical queries
- Efficient at counting distinct values
- Good support for time-series data
- Handles large data volumes efficiently
- Supports materialized views for faster query performance
- **AggregatingMergeTree**: Enhances materialized views by pre-aggregating data, reducing query time and improving performance for large datasets

#### 4.3 Why Batch Processing?

- Reduces database write load
- Improves throughput
- More efficient resource usage
- Better error handling

#### 4.5 Why gRPC?

- **Performance**: Better performance compared to REST due to:
  - Binary protocol (Protocol Buffers)
  - HTTP/2 multiplexing
  - Stream support
- **Strong Typing**: Protocol Buffers provide type safety and validation
- **Bi-directional Streaming**: Useful for real-time segment updates
- **Code Generation**: Automatic client/server code generation reduces errors
- **Service Definition**: Clear contract between services using .proto files

### 5. Scalability Considerations

#### 5.1 Horizontal Scaling

- USS can scale independently
- Multiple ES instances possible
- Kafka partitioning for parallelism
- ClickHouse distributed capabilities

#### 5.2 Performance Optimization

- Batch processing reduces load
- Efficient storage in ClickHouse
- Kafka as buffer for traffic spikes
- Index optimization in ClickHouse

### 6. Error Handling

#### 6.1 Message Processing

- Retry mechanism for failed batches
- Dead letter queue for invalid messages
- Error logging and monitoring
- Graceful degradation

#### 6.2 Data Consistency

- At-least-once delivery by using Group in Kafka
- Idempotent processing
- Transaction boundaries
- Data validation

### 7. Configuration Parameters

This service is designed following microservices principles, allowing it to be easily configured. The configuration parameters (such as environment variables) are defined in the following locations:

Path: Configuration variables are found in the config package.
Configuration Files:

- env.go: Defines general environment variables.
- kafka_env.go: Defines Kafka-related configuration parameters.
  These variables are loaded at runtime, allowing you to adjust the service's behavior by modifying the configuration files or the respective environment variables.

Example environment variables:

SERVICE_HOST: Host address of the service
SERVICE_PORT: Port the service listens on
KAFKA_BROKER_URL: Kafka broker address
KAFKA_TOPIC: Kafka topic for segment messages
CLICKHOUSE_ADDR: ClickHouse database URL

### 8. Testing Strategy

Unit Tests:

- Business logic
- Data transformations
- Error handling

### 9. Makefile Commands

You can manage the service's migrations, versioning, tests, and linting using the Makefile. Here are some commands available:

- migrate: Run make migrate to apply database migrations. This ensures that all schema changes are applied to your ClickHouse database.
  Versioning: Run make version to check the current version of the service.
- test: Use make test to run all unit and integration tests, ensuring that the service functions as expected.
- lint: Use make lint to run the code linter and ensure that the code adheres to best practices and style guidelines.
  These commands allow for easy and consistent management of the project's lifecycle and ensure that the system is always in a working state.

### 9. Future Improvements

1. Monitoring and Metrics

   - Batch processing stats
   - Processing latency
   - Error rates
   - Resource usage

2. Enhanced Features

   - Segment analytics
   - User behavior tracking
   - Real-time notifications

3. User Authentication:

   - User identity verification
   - Role-based access control (RBAC)
   - Session management

4. Integration Tests

- End-to-end data flow
- Service interactions
- Kafka message processing

5. Benchmark Tests

- Performance under load
- Throughput and latency
- Resource utilization
- Scalability testing
