// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import "time"

// MetricsSink receives collected metrics for external processing.
type MetricsSink func(sinkCtx SinkContext, metrics Metrics) error

// SinkContext carries instance-level identity for metrics aggregation.
type SinkContext struct {
	TopicName       string
	ConsumerGroup   string
	ApplicationName string
	Team            *Team // nil if not configured via WithTeam()
}

// Metrics mutated during message processing, with Traits and other
// envelope information copied prior to MetricsSink call.
//
// COMPLIANCE NOTE: metrics are for observability only, partition keys
// are deliberately excluded to prevent accidental sensitive data disclosure
type Metrics struct {
	Traits                  Traits // framework may add others, e.g. DeadLetter, ProcessError...
	QueueDepth              int32  // for buffering implementations; int32 to fit struct in 2 lines
	Partition               int32
	Offset                  int64
	ProcessDuration         time.Duration
	WriteDeadLetterDuration time.Duration
	ProcessStartTime        time.Time
	ReadTime                time.Time
	WatermarkAdvanceTime    time.Time
	WorkerPool              uint32
}

// AddCustomTraits which are merged with Message traits
// prior to metrics collection - idempotent
func (m *Metrics) AddCustomTraits(customTraits Traits) {
	m.Traits |= customTraits &^ FrameworkReserved
}
