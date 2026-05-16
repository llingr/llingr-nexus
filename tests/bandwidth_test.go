// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/llingr/llingr-nexus/nexus"
)

func TestBandwidthMetricsSink(t *testing.T) {
	t.Run("invocation with populated metrics", func(t *testing.T) {
		var receivedTopic string
		var receivedMetrics nexus.BandwidthMetrics

		sink := nexus.BandwidthMetricsSink(func(topicName string, metrics nexus.BandwidthMetrics) error {
			receivedTopic = topicName
			receivedMetrics = metrics
			return nil
		})

		sent := nexus.BandwidthMetrics{
			BandwidthMetricsID: "test-uuid-123",
			TopicName:          "orders",
			ConsumerGroup:      "order-processor",
		}

		err := sink("orders", sent)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if receivedTopic != "orders" {
			t.Errorf("expected topic 'orders', got %q", receivedTopic)
		}
		if receivedMetrics.BandwidthMetricsID != "test-uuid-123" {
			t.Errorf("expected BandwidthMetricsID 'test-uuid-123', got %q", receivedMetrics.BandwidthMetricsID)
		}
		if receivedMetrics.ConsumerGroup != "order-processor" {
			t.Errorf("expected ConsumerGroup 'order-processor', got %q", receivedMetrics.ConsumerGroup)
		}
	})

	t.Run("propagates error from sink", func(t *testing.T) {
		sinkErr := errors.New("sink unavailable")
		sink := nexus.BandwidthMetricsSink(func(_ string, _ nexus.BandwidthMetrics) error {
			return sinkErr
		})

		err := sink("topic", nexus.BandwidthMetrics{})
		if !errors.Is(err, sinkErr) {
			t.Errorf("expected sink error, got %v", err)
		}
	})

	t.Run("receives zero-value metrics without panic", func(t *testing.T) {
		called := false
		sink := nexus.BandwidthMetricsSink(func(_ string, _ nexus.BandwidthMetrics) error {
			called = true
			return nil
		})

		err := sink("", nexus.BandwidthMetrics{})
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if !called {
			t.Error("sink should have been called")
		}
	})
}

func TestBandwidthMetrics(t *testing.T) {
	t.Run("zero value is valid", func(t *testing.T) {
		var m nexus.BandwidthMetrics
		if m.Ts != (time.Time{}) {
			t.Error("zero Ts should be zero time")
		}
		if m.StatsIntervalDuration != 0 {
			t.Error("zero StatsIntervalDuration should be 0")
		}
		if m.BandwidthMetricsID != "" {
			t.Error("zero BandwidthMetricsID should be empty")
		}
		if m.TopicName != "" {
			t.Error("zero TopicName should be empty")
		}
		if m.ConsumerGroup != "" {
			t.Error("zero ConsumerGroup should be empty")
		}
		if m.Brokers != nil {
			t.Error("zero Brokers should be nil")
		}
		if m.Partitions != nil {
			t.Error("zero Partitions should be nil")
		}
	})

	t.Run("fields are accessible after construction", func(t *testing.T) {
		now := time.Now()
		m := nexus.BandwidthMetrics{
			Ts:                    now,
			StatsIntervalDuration: time.Minute,
			BandwidthMetricsID:    "abc-123",
			TopicName:             "events",
			ConsumerGroup:         "consumer-1",
			Brokers: []nexus.BrokerInfo{
				{ID: "1", Host: "broker-1", Port: "9092", Rack: "us-east-1a"},
			},
			Partitions: []nexus.PartitionBandwidth{
				{ID: 0, ReceivedBytes: 1024},
			},
		}

		if !m.Ts.Equal(now) {
			t.Error("Ts mismatch")
		}
		if m.StatsIntervalDuration != time.Minute {
			t.Error("StatsIntervalDuration mismatch")
		}
		if m.BandwidthMetricsID != "abc-123" {
			t.Error("BandwidthMetricsID mismatch")
		}
		if m.TopicName != "events" {
			t.Error("TopicName mismatch")
		}
		if m.ConsumerGroup != "consumer-1" {
			t.Error("ConsumerGroup mismatch")
		}
		if len(m.Brokers) != 1 {
			t.Fatalf("expected 1 broker, got %d", len(m.Brokers))
		}
		if len(m.Partitions) != 1 {
			t.Fatalf("expected 1 partition, got %d", len(m.Partitions))
		}
	})

	t.Run("with multiple brokers and partitions", func(t *testing.T) {
		m := nexus.BandwidthMetrics{
			Brokers: []nexus.BrokerInfo{
				{ID: "1", Host: "broker-1", Port: "9092", Rack: "us-east-1a"},
				{ID: "2", Host: "broker-2", Port: "9092", Rack: "us-east-1b"},
				{ID: "3", Host: "broker-3", Port: "9092", Rack: "us-east-1c"},
			},
			Partitions: []nexus.PartitionBandwidth{
				{ID: 0, ReceivedBytes: 1000},
				{ID: 1, ReceivedBytes: 2000},
				{ID: 2, ReceivedBytes: 3000},
				{ID: 3, ReceivedBytes: 4000},
			},
		}

		if len(m.Brokers) != 3 {
			t.Errorf("expected 3 brokers, got %d", len(m.Brokers))
		}
		if len(m.Partitions) != 4 {
			t.Errorf("expected 4 partitions, got %d", len(m.Partitions))
		}
	})

	t.Run("empty brokers and partitions", func(t *testing.T) {
		m := nexus.BandwidthMetrics{
			Brokers:    []nexus.BrokerInfo{},
			Partitions: []nexus.PartitionBandwidth{},
		}

		if len(m.Brokers) != 0 {
			t.Error("expected empty brokers slice")
		}
		if len(m.Partitions) != 0 {
			t.Error("expected empty partitions slice")
		}
	})
}

func TestBrokerInfo(t *testing.T) {
	t.Run("zero value is valid", func(t *testing.T) {
		var b nexus.BrokerInfo
		if b.ID != "" || b.Host != "" || b.Port != "" || b.Rack != "" {
			t.Error("zero BrokerInfo should have all empty strings")
		}
	})

	t.Run("fields are accessible", func(t *testing.T) {
		b := nexus.BrokerInfo{
			ID:   "101",
			Host: "kafka-broker-101.internal",
			Port: "9092",
			Rack: "eu-west-2a",
		}

		if b.ID != "101" {
			t.Errorf("expected ID '101', got %q", b.ID)
		}
		if b.Host != "kafka-broker-101.internal" {
			t.Errorf("expected Host 'kafka-broker-101.internal', got %q", b.Host)
		}
		if b.Port != "9092" {
			t.Errorf("expected Port '9092', got %q", b.Port)
		}
		if b.Rack != "eu-west-2a" {
			t.Errorf("expected Rack 'eu-west-2a', got %q", b.Rack)
		}
	})

	t.Run("rack may be empty", func(t *testing.T) {
		b := nexus.BrokerInfo{
			ID:   "1",
			Host: "localhost",
			Port: "9092",
			Rack: "",
		}

		if b.Rack != "" {
			t.Error("empty Rack should be valid")
		}
	})
}

func TestPartitionBandwidth(t *testing.T) {
	t.Run("zero value is valid", func(t *testing.T) {
		var p nexus.PartitionBandwidth
		if p.Ts != (time.Time{}) {
			t.Error("zero Ts should be zero time")
		}
		if p.ReceivedBytes != 0 || p.TransmittedBytes != 0 || p.ReceivedMessageCount != 0 {
			t.Error("zero counters should be 0")
		}
		if p.CompressedBytes != 0 || p.UncompressedBytes != 0 {
			t.Error("zero compression bytes should be 0")
		}
		if p.ID != 0 {
			t.Error("zero ID should be 0")
		}
		if p.Leader != "" || p.Compression != "" {
			t.Error("zero strings should be empty")
		}
	})

	t.Run("fields are accessible after construction", func(t *testing.T) {
		now := time.Now()
		p := nexus.PartitionBandwidth{
			Ts:                   now,
			ReceivedBytes:        1048576,
			TransmittedBytes:     524288,
			ReceivedMessageCount: 10000,
			CompressedBytes:      262144,
			UncompressedBytes:    1048576,
			ID:                   7,
			Leader:               "3",
			Compression:          "snappy",
		}

		if !p.Ts.Equal(now) {
			t.Error("Ts mismatch")
		}
		if p.ReceivedBytes != 1048576 {
			t.Error("ReceivedBytes mismatch")
		}
		if p.TransmittedBytes != 524288 {
			t.Error("TransmittedBytes mismatch")
		}
		if p.ReceivedMessageCount != 10000 {
			t.Error("ReceivedMessageCount mismatch")
		}
		if p.CompressedBytes != 262144 {
			t.Error("CompressedBytes mismatch")
		}
		if p.UncompressedBytes != 1048576 {
			t.Error("UncompressedBytes mismatch")
		}
		if p.ID != 7 {
			t.Error("ID mismatch")
		}
		if p.Leader != "3" {
			t.Error("Leader mismatch")
		}
		if p.Compression != "snappy" {
			t.Error("Compression mismatch")
		}
	})

	t.Run("compression fields zero when unavailable", func(t *testing.T) {
		p := nexus.PartitionBandwidth{
			ID:                   0,
			ReceivedBytes:        4096,
			TransmittedBytes:     2048,
			ReceivedMessageCount: 50,
			CompressedBytes:      0,
			UncompressedBytes:    0,
			Compression:          "",
		}

		if p.CompressedBytes != 0 {
			t.Error("CompressedBytes should be zero when unavailable")
		}
		if p.UncompressedBytes != 0 {
			t.Error("UncompressedBytes should be zero when unavailable")
		}
		if p.Compression != "" {
			t.Error("Compression should be empty when unavailable")
		}
	})

	t.Run("partition ID matches int32 range", func(t *testing.T) {
		p := nexus.PartitionBandwidth{ID: math.MaxInt32}
		if p.ID != math.MaxInt32 {
			t.Errorf("expected max int32 (%d), got %d", int32(math.MaxInt32), p.ID)
		}
	})

	t.Run("cumulative counters are int64", func(t *testing.T) {
		p := nexus.PartitionBandwidth{
			ReceivedBytes:        math.MaxInt64,
			TransmittedBytes:     math.MaxInt64,
			ReceivedMessageCount: math.MaxInt64,
		}

		if p.ReceivedBytes != math.MaxInt64 {
			t.Error("ReceivedBytes should support full int64 range")
		}
		if p.TransmittedBytes != math.MaxInt64 {
			t.Error("TransmittedBytes should support full int64 range")
		}
		if p.ReceivedMessageCount != math.MaxInt64 {
			t.Error("ReceivedMessageCount should support full int64 range")
		}
	})
}

func TestValidateBandwidthInterval(t *testing.T) {
	t.Run("zero duration returns nil", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(0); err != nil {
			t.Errorf("zero should return nil (use default), got %v", err)
		}
	})

	t.Run("one second is valid", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(time.Second); err != nil {
			t.Errorf("1s should be valid, got %v", err)
		}
	})

	t.Run("one minute is valid", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(time.Minute); err != nil {
			t.Errorf("1m should be valid, got %v", err)
		}
	})

	t.Run("twelve hours is valid", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(12 * time.Hour); err != nil {
			t.Errorf("12h should be valid, got %v", err)
		}
	})

	t.Run("sub-second returns error", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(999 * time.Millisecond); err == nil {
			t.Error("999ms should return error")
		}
		if err := nexus.ValidateBandwidthInterval(time.Nanosecond); err == nil {
			t.Error("1ns should return error")
		}
	})

	t.Run("exceeds twelve hours returns error", func(t *testing.T) {
		if err := nexus.ValidateBandwidthInterval(12*time.Hour + time.Nanosecond); err == nil {
			t.Error("12h + 1ns should return error")
		}
	})

	t.Run("negative duration returns error", func(t *testing.T) {
		err := nexus.ValidateBandwidthInterval(-time.Second)
		if err == nil {
			t.Error("negative duration should return error")
		}
		if err.Error() != "bandwidth interval must not be negative" {
			t.Errorf("expected negative-specific error message, got %q", err.Error())
		}
	})
}

func TestDefaultBandwidthInterval_Value(t *testing.T) {
	if nexus.DefaultBandwidthInterval != time.Minute {
		t.Errorf("DefaultBandwidthInterval should be 1 minute, got %v", nexus.DefaultBandwidthInterval)
	}
}

// mockBandwidthAdapter is a compile-time interface satisfaction check
// and test double for BandwidthPort[string].
type mockBandwidthAdapter struct {
	callback nexus.BandwidthCallback
	interval time.Duration
}

var _ nexus.BandwidthPort[string] = (*mockBandwidthAdapter)(nil)

func (m *mockBandwidthAdapter) SetBandwidthCallback(cb nexus.BandwidthCallback) {
	m.callback = cb
}

func (m *mockBandwidthAdapter) StatsInterval() time.Duration {
	return m.interval
}

func TestBandwidthPort(t *testing.T) {
	t.Run("SetBandwidthCallback receives data", func(t *testing.T) {
		adapter := &mockBandwidthAdapter{interval: time.Minute}

		var received nexus.BandwidthMetrics
		adapter.SetBandwidthCallback(func(m nexus.BandwidthMetrics) {
			received = m
		})

		sent := nexus.BandwidthMetrics{
			BandwidthMetricsID: "callback-test",
			TopicName:          "events",
		}
		adapter.callback(sent)

		if received.BandwidthMetricsID != "callback-test" {
			t.Errorf("expected BandwidthMetricsID 'callback-test', got %q", received.BandwidthMetricsID)
		}
		if received.TopicName != "events" {
			t.Errorf("expected TopicName 'events', got %q", received.TopicName)
		}
	})

	t.Run("StatsInterval returns configured duration", func(t *testing.T) {
		adapter := &mockBandwidthAdapter{interval: 30 * time.Second}

		if adapter.StatsInterval() != 30*time.Second {
			t.Errorf("expected 30s, got %v", adapter.StatsInterval())
		}
	})
}
