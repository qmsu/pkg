package utils

import "testing"

func TestRandStringRunes(t *testing.T) {
	randNumMap := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		randNum := RandStringRunes(8, false)
		if _, ok := randNumMap[randNum]; ok {
			t.Fatal("failed:", randNum)
		}
	}
}
