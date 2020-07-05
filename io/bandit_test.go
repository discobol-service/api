package io

import (
	"testing"
)

func TestGetBandit (t *testing.T) {
	bandit := GetBandit("http://localhost:8888")

	if bandit.Address == "" {
		t.Error("Address not found for bandit")
	}
}
