// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

//go:embed service_test.json
var serviceJSON []byte

// exampleServiceSpec is a Service.Spec payload as an application might define it.
// Nexus passes the JSON bytes through unchanged.
type exampleServiceSpec struct {
	Description   string              `json:"description,omitempty"`
	Lifecycle     string              `json:"lifecycle,omitempty"`
	Language      string              `json:"language,omitempty"`
	Repository    string              `json:"repository,omitempty"`
	Documentation string              `json:"documentation,omitempty"`
	Runbook       string              `json:"runbook,omitempty"`
	Maintainers   []string            `json:"maintainers,omitempty"`
	Oncall        exampleOncall       `json:"oncall,omitempty"`
	Contacts      []exampleContact    `json:"contacts,omitempty"`
	Dependencies  exampleDependencies `json:"dependencies,omitempty"`
	SLA           exampleSLA          `json:"sla,omitempty"`
	Tags          map[string]string   `json:"tags,omitempty"`
}

type exampleOncall struct {
	Primary    string `json:"primary"`
	Escalation string `json:"escalation,omitempty"`
	Schedule   string `json:"schedule,omitempty"`
}

type exampleContact struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Address string `json:"address"`
}

type exampleDependencies struct {
	Upstream   []string `json:"upstream,omitempty"`
	Downstream []string `json:"downstream,omitempty"`
}

type exampleSLA struct {
	Availability float64 `json:"availability"`
	LatencyP99Ms int     `json:"latency_p99_ms"`
}

// TestService_DecodeFromJSON shows the end-to-end pattern for consumers
// loading a Service from JSON: unmarshal the outer Service with encoding/json,
// then unmarshal the inner Spec bytes into the application's own type
func TestService_DecodeFromJSON(t *testing.T) {
	var svc nexus.Service
	if err := json.Unmarshal(serviceJSON, &svc); err != nil {
		t.Fatalf("unmarshal service JSON: %v", err)
	}

	if svc.Name != "payments-api" {
		t.Errorf("Name = %q, want payments-api", svc.Name)
	}
	if svc.Team != "payments-team" {
		t.Errorf("Team = %q, want payments-team", svc.Team)
	}

	var spec exampleServiceSpec
	if err := json.Unmarshal(svc.Spec, &spec); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}

	if spec.Lifecycle != "production" {
		t.Errorf("Lifecycle = %q, want production", spec.Lifecycle)
	}
	if got, want := len(spec.Maintainers), 2; got != want {
		t.Errorf("len(Maintainers) = %d, want %d", got, want)
	}
	if spec.Oncall.Primary != "payments-oncall" {
		t.Errorf("Oncall.Primary = %q, want payments-oncall", spec.Oncall.Primary)
	}
	if got, want := len(spec.Contacts), 4; got != want {
		t.Fatalf("len(Contacts) = %d, want %d", got, want)
	}
	if spec.Contacts[1].Type != "slack" || spec.Contacts[1].Label != "alerts" {
		t.Errorf("Contacts[1] = %+v, want {slack, alerts, ...}", spec.Contacts[1])
	}
	if got, want := len(spec.Dependencies.Downstream), 2; got != want {
		t.Errorf("len(Dependencies.Downstream) = %d, want %d", got, want)
	}
	if spec.SLA.Availability != 99.95 {
		t.Errorf("SLA.Availability = %v, want 99.95", spec.SLA.Availability)
	}
	if spec.Tags["domain"] != "commerce" {
		t.Errorf("Tags[domain] = %q, want commerce", spec.Tags["domain"])
	}
}
