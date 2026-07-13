# llingr-nexus

[![CI](https://github.com/llingr/llingr-nexus/actions/workflows/ci.yml/badge.svg)](https://github.com/llingr/llingr-nexus/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/llingr/llingr-nexus.svg)](https://pkg.go.dev/github.com/llingr/llingr-nexus)
[![Tag](https://img.shields.io/github/v/tag/llingr/llingr-nexus)](https://github.com/llingr/llingr-nexus/tags)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/llingr/llingr-nexus)](go.mod)

Shared contracts for broker-agnostic message processing.

## Overview

llingr-nexus provides a clean abstraction layer that can be implemented by message
processing engines, and used by broker-specific adapters. Applications are presented
with stable interfaces that work across message brokers with partition/offset
semantics.

The processing engine delivers messages to the application via `ProcessMessage[T](ctx, msg)`.
Any engine implementing these contracts works with any adapter implementing `BrokerPort`.

## Installation

```bash
go get github.com/llingr/llingr-nexus
```

## Port-Binding Builder Pattern

llingr-nexus implements a Port-Binding Builder pattern - a dependency
injection approach inspired by Hexagonal Architecture (Alistair Cockburn, 2005) where
the builder carries application dependencies to an adapter, which then injects itself
to complete construction.

```go
// 1. Application creates a builder with its callbacks and configuration
builder := engine.NewBuilder[BrokerMessage](processMessage, writeDeadLetter).
    WithConfig(cfg).
    WithMetricsSink(metricsSink)

// 2. Adapter receives builder, calls Build(self), returns wired consumer
adapter := brokeradapter.New(ctx, "my-group", "broker:9092")
consumer, _ := adapter.CreateConsumer(builder)

// 3. Subscribe and run
consumer.Subscribe()
```

The pattern enables loose coupling between processing engines and adapters - they
communicate only through shared contracts, making components independently testable,
replaceable, and composable. Wiring is verified at compile time.

## Core Interfaces

### BrokerPort

The adapter interface. Implementations connect the consumer to specific message brokers
(Kafka, Pulsar, NATS JetStream, or any system with partition/offset semantics).

```go
type BrokerPort[T any] interface {
    Subscribe() error
    Unsubscribe() error
    Poll(timeout time.Duration) (T, bool, error)
    ExtractEnvelope(T) Envelope
    CommitOffsets(messages []*Message[T]) ([]*Message[T], error)
    AckRebalance(RebalanceType, []RebalanceInfo) error
    BrokerQuery(QueryRequest) (QueryResponse, error)
}
```

`BrokerQuery` is a generic broker-side query mechanism. The only `QueryType`
currently defined is `CommittedOffsets`; additional types can be introduced
without breaking the adapter contract.

### ConsumerBuilder

Carries host application dependencies across module boundaries. The adapter calls
`Build()` to inject itself as `BrokerPort`, completing construction.

```go
type ConsumerBuilder[T any] interface {
    Build(BrokerPort[T]) AdaptedConsumer[T]
    TopicName() string
}
```

`Build()` is the binding site - compiler-verified type safety ensures correct wiring.
`TopicName()` exposes the configured topic to adapters.

### Consumer

The public interface returned to the host application for lifecycle control.

```go
type Consumer[T any] interface {
    Subscribe() error
    Shutdown() error
}
```

## Core Types

### Message

Broker-agnostic message with addressing info and optional business intelligence traits.
Uses `*T` for sync.Pool compatibility.

```go
type Message[T any] struct {
    Traits    Traits // bit flags, positions 0-9 reserved
    Partition int32
    Offset    int64
    Key       string
    Payload   *T     // broker-specific message type
}
```

### Envelope

Extracted addressing information from broker-specific payload.

```go
type Envelope struct {
    Partition int32
    Offset    int64
    Key       string
    Ctx       context.Context // for traces/spans
}
```

### Metrics

Self-contained observability data. Fields aligned for efficient packing
within the minimum number of cache lines. Keys are deliberately excluded
to prevent accidental PII disclosure.

```go
type Metrics struct {
    Traits                  Traits
    QueueDepth              int32
    Partition               int32
    Offset                  int64
    ProcessDuration         time.Duration
    WriteDeadLetterDuration time.Duration
    ReadTime                time.Time
    ProcessStartTime        time.Time
    WatermarkAdvanceTime    time.Time
}
```

## Application Callbacks

Two functions are required when creating a consumer builder:

### ProcessMessage

Blocking call to process each message. Must not spawn goroutines that outlive the call.

```go
type ProcessMessage[T any] func(ctx context.Context, msg *Message[T]) error
```

### WriteDeadLetter

Handles failed messages. Called when ProcessMessage returns an error.

```go
type WriteDeadLetter[T any] func(ctx context.Context, msg *Message[T], reason error) error
```

### MetricsSink (optional)

Receives metrics for each processed message. Fire-and-forget pattern.

```go
type MetricsSink func(topicName string, metrics Metrics) error
```

## Traits

64-bit flags for processing engine state and custom business intelligence.

### Reserved Traits (bits 0-9)

| Bit | Trait               | Description                            |
|-----|---------------------|----------------------------------------|
| 0   | ProcessError        | Processing returned an error           |
| 1   | ProcessPanic        | Processing panicked                    |
| 2   | DeadLetter          | Message sent to dead letter handler    |
| 3   | CommitBuffered      | Commit buffered in gap buffer          |
| 4   | Duplicate           | Duplicate message detected             |
| 5   | UsedOverflow        | Used shared overflow capacity          |
| 6   | Orphaned            | Completed after partition reassignment |
| 7   | FirstAfterRebalance | First message after rebalance          |
| 8   | _reserved_          | Reserved for future framework use      |
| 9   | _reserved_          | Reserved for future framework use      |

### Custom Traits (bits 10-63)

Applications can define custom traits for business intelligence:

```go
const (
    HighValue   nexus.Traits = 1 << 10  // custom trait
    FraudAlert  nexus.Traits = 1 << 11  // custom trait
)

// Add to message during processing
msg.AddCustomTraits(HighValue | FraudAlert)
```

### Checking Traits

```go
// Go - single trait
if metrics.Traits&nexus.DeadLetter != 0 {
    // handle dead letter
}

// Go - multiple traits (any match)
const alertCondition = nexus.DeadLetter | nexus.ProcessPanic
if metrics.Traits&alertCondition != 0 {
    // at least one alert condition present
}
```

```python
# Python - in analytics pipelines
DEAD_LETTER = 1 << 2
FRAUD_ALERT = 1 << 11

df['needs_review'] = df['traits'].apply(
    lambda t: (t & (DEAD_LETTER | FRAUD_ALERT)) != 0
)
```

## Design Principles

- **Zero dependencies** - only Go stdlib
- **Type safety** - generics throughout, compile-time guarantees
- **Performance** - cache-line aligned structs, sync.Pool compatible
- **Privacy** - keys excluded from metrics by design
- **Stability** - interface changes are rare and deliberate
- **Broker agnostic** - works with any partition/offset message system

## Licence

Apache-2.0 - see [LICENSE](./LICENSE) and [COPYRIGHT](./COPYRIGHT).
Contributions are governed by [CONTRIBUTING.md](./CONTRIBUTING.md).
