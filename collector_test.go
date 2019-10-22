package main

import "testing"

func TestBoolToFloat64(t *testing.T) {
	trueResult := boolToFloat64(true)
	if trueResult != 1 {
		t.Errorf("boolToFloat64 was incorrect. expected %v got %v", 1, trueResult)
	}

	falseResult := boolToFloat64(false)
	if falseResult != 0 {
		t.Errorf("boolToFloat64 was incorrect. expected %v got %v", 0, falseResult)
	}

}
