// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import (
	"encoding/json"
)

// Service identifies a microservice or consumer instance,
// with optional ownership and declarative metadata.
type Service struct {
	Name string          `json:"name"`           // typically the application name for the consumer
	Team string          `json:"team,omitempty"` // the primary/responsible owner of the Service
	Spec json.RawMessage `json:"spec,omitempty"`
}
