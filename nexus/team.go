// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

// Team identifies the owning team for a consumer instance.
// Declared at build time via WithTeam() and carried through
// the metrics pipeline to fleet discovery.
type Team struct {
	Name       string    `json:"name"`
	Department string    `json:"department,omitempty"`
	Channels   []Channel `json:"channels"`
}

// Channel is a communication endpoint for a team.
type Channel struct {
	Label    string      `json:"label"`
	Address  string      `json:"address"`
	Type     ChannelType `json:"type"`
	Purposes []Purpose   `json:"purposes"`
}

// ChannelType identifies the communication platform.
type ChannelType string

const (
	Slack      ChannelType = "slack"
	PagerDuty  ChannelType = "pagerduty"
	Zendesk    ChannelType = "zendesk"
	OpsGenie   ChannelType = "opsgenie"
	MsTeams    ChannelType = "teams"
	Discord    ChannelType = "discord"
	GoogleChat ChannelType = "googlechat"
	Email      ChannelType = "email"
	Webhook    ChannelType = "webhook"
)

// Purpose describes what a channel is used for.
type Purpose string

const (
	Primary    Purpose = "primary"
	Alerts     Purpose = "alerts"
	Escalation Purpose = "escalation"
	Reports    Purpose = "reports"
)
