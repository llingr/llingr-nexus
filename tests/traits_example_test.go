// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests_test

import (
	"fmt"

	"github.com/llingr/llingr-nexus/nexus"
)

// some trait examples - illustrative only, your use-cases will vary
const (
	vipCustomer nexus.Traits = 1 << (10 + iota)
)

const (
	euCustomer nexus.Traits = 1 << (20 + iota)
)

const (
	priorityProcessing nexus.Traits = 1 << (30 + iota)
	auditRequired
	streamWorthy
)

const (
	highValueTransaction nexus.Traits = 1 << (40 + iota)
	fraudSuspicion
)

func ExampleSetTraits_basicUsage() {
	var traits nexus.Traits

	nexus.SetTraits(&traits, euCustomer)
	if traits&euCustomer != 0 {
		fmt.Println("customer requires GDPR compliance")
	}

	// Output:
	// customer requires GDPR compliance
}

func ExampleSetTraits_efficientMultipleFlags() {
	var traits nexus.Traits

	// efficient multiple flag setting in one operation
	processingFlags := auditRequired | priorityProcessing
	nexus.SetTraits(&traits, processingFlags)

	if traits&auditRequired != 0 {
		fmt.Println("audit enabled")
	}
	if traits&priorityProcessing != 0 {
		fmt.Println("priority processing enabled")
	}

	// Output:
	// audit enabled
	// priority processing enabled
}

func ExampleSetTraits_alertConditions() {
	var traits nexus.Traits

	// set up a critical scenario: VIP customer message failed
	nexus.SetDeadLetter(&traits)
	nexus.SetTraits(&traits, vipCustomer)
	nexus.SetTraits(&traits, highValueTransaction)

	// illustrative alert conditions
	criticalAlert := nexus.DeadLetter | vipCustomer
	highValueAlert := highValueTransaction | vipCustomer

	if traits&criticalAlert == criticalAlert {
		fmt.Println("critical alert: vip customer message failed")
	}

	if traits&highValueAlert == highValueAlert {
		fmt.Println("high value vip transaction detected")
	}

	// Output:
	// critical alert: vip customer message failed
	// high value vip transaction detected
}

func ExampleUnsetTraits_clearSensitiveFlags() {
	var traits nexus.Traits

	// transaction initially flagged with multiple concerns
	nexus.SetTraits(&traits, vipCustomer|highValueTransaction)
	nexus.SetTraits(&traits, fraudSuspicion|streamWorthy)

	fmt.Println("before investigation: fraud flags present")

	// after investigation completes, clear sensitive flags but keep business traits
	sensitiveFlags := fraudSuspicion | streamWorthy
	nexus.UnsetTraits(&traits, sensitiveFlags)

	// verify sensitive flags are cleared but business traits remain
	if traits&fraudSuspicion == 0 && traits&streamWorthy == 0 {
		fmt.Println("investigation complete: sensitive flags cleared")
	}
	if traits&vipCustomer != 0 && traits&highValueTransaction != 0 {
		fmt.Println("business traits preserved for future processing")
	}

	// Output:
	// before investigation: fraud flags present
	// investigation complete: sensitive flags cleared
	// business traits preserved for future processing
}
