// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

// RebalanceType currently includes assign and revoke.
// In future may add LostPartitions and others as needed.
type RebalanceType int

// RebalanceType values indicating partition assignment changes.
const (
	_      RebalanceType = iota // 0 - reserved
	Assign                      // 1 - partitions assigned
	Revoke                      // 2 - partitions revoked
)

// RebalanceInfo for a specific topic/partition
// with optional broker-specific metadata.
type RebalanceInfo struct {
	RebalanceType   RebalanceType
	TopicName       string
	Partition       int32
	CommittedOffset int64 // optional, offset of **next** record to be processed
	Meta            any   // optional adapter-specific information cast to concrete type as needed
}
