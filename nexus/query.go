// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

// QueryType for broker QueryRequest / QueryResponse
type QueryType int

// CommittedOffsets query type, currently
// the only supported query
const (
	_                QueryType = iota // 0
	CommittedOffsets                  // 1
)

// QueryRequest for broker queries
type QueryRequest struct {
	QueryType       QueryType // current implementation only supports CommittedOffsets
	TopicName       string    //
	TopicPartitions []int32   // to limit query scope, empty typically considers assigned partitions
	Data            any       // adapter-specific, adapter casts to concrete type
}

// QueryResponse from broker queries
type QueryResponse struct {
	QueryType          QueryType       // current implementation only supports CommittedOffsets
	TopicName          string          //
	OffsetsByPartition map[int32]int64 // for CommittedOffsets queries
	Data               any             // adapter-specific, adapter casts to concrete type
}
