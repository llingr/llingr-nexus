// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import (
	"errors"
	"time"
)

// DefaultBandwidthInterval metrics collection cadence.
const DefaultBandwidthInterval = time.Minute

// BandwidthMetricsSink receives telemetry for external processing.
type BandwidthMetricsSink func(topicName string, metrics BandwidthMetrics) error

// BandwidthMetrics carries per-partition cumulative bandwidth byte counts
// at a given measurement time.
//
// COMPLIANCE NOTE: bandwidth metrics contain only infrastructure-level
// information (bytes, message counts, and addressable broker topology).
type BandwidthMetrics struct {
	Ts                    time.Time     // measurement timestamp
	StatsIntervalDuration time.Duration // adapter collection cadence
	BandwidthMetricsID    string        // idempotency ID (typically UUID)
	TopicName             string
	ConsumerGroup         string
	Brokers               []BrokerInfo
	Partitions            []PartitionBandwidth
	Team                  *Team // nil if not configured via WithTeam()
}

// BrokerInfo describes a broker node at the time of measurement.
// All identifiers are strings for forward compatibility — Kafka
// uses int32 for broker IDs, but other systems may not.
type BrokerInfo struct {
	ID   string // broker node ID
	Host string
	Port string
	Rack string // availability zone or rack; empty when not supported in broker/adapter
}

// PartitionBandwidth carries cumulative bandwidth byte counts for a single partition
// over one collection interval.
//
// CompressedBytes, UncompressedBytes, and Compression can be provided by client adapters
// that expose compression visibility.
// When unavailable these fields will be zero/empty
type PartitionBandwidth struct {
	Ts                   time.Time // measurement timestamp
	ReceivedBytes        int64     // cumulative bytes received
	TransmittedBytes     int64     // cumulative bytes transmitted
	ReceivedMessageCount int64     // cumulative messages received
	CompressedBytes      int64     // wire bytes
	UncompressedBytes    int64     // decompressed bytes
	ID                   int32     // partition identifier - aligns with Metrics.Partition
	Leader               string    // broker serving this partition
	Compression          string    // algorithm name
}

// BandwidthCallback receives telemetry for delivery to BandwidthMetricsSink
type BandwidthCallback func(BandwidthMetrics)

// BandwidthPort is an optional interface that adapters may implement
// alongside BrokerPort[T] to supply bandwidth telemetry if a given
// broker client/adapter supports this.
//
// The adapter is the pump: it controls the collection cadence and invokes
// the callback at each interval tick. The consumer builder registers the
// callback and routes data to the user's BandwidthMetricsSink.
type BandwidthPort[T any] interface {
	// SetBandwidthCallback registers the function the adapter calls on
	// each stats interval tick. The consumer builder provides this
	// callback - adapters must not call it before registration.
	SetBandwidthCallback(BandwidthCallback)

	// StatsInterval returns the adapter's configured collection cadence.
	StatsInterval() time.Duration
}

// ValidateBandwidthInterval checks that d is within the allowed range
// for bandwidth collection cadence. Returns nil for zero (meaning "use
// DefaultBandwidthInterval"). Valid range: [1s, 12h].
func ValidateBandwidthInterval(d time.Duration) error {
	if d == 0 {
		return nil // zero means "use default"
	}
	if d < 0 {
		return errors.New("bandwidth interval must not be negative")
	}
	if d < time.Second {
		return errors.New("bandwidth interval must be at least 1 second")
	}
	if d > 12*time.Hour {
		return errors.New("bandwidth interval must not exceed 12 hours")
	}
	return nil
}
