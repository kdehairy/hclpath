package cmpval

import (
	"fmt"
	"math"
	"strconv"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func cmpString(val cty.Value, expected string) bool {
	return expected == val.AsString()
}

func cmpNumber(val cty.Value, expected float64) (bool, error) {
	const tolerance = 1e-9
	var a float64
	err := gocty.FromCtyValue(val, &a)
	if err != nil {
		return false, fmt.Errorf("failed to parse number: %v", err)
	}
	return math.Abs(expected-a) <= tolerance, nil
}

func IsEqual(val cty.Value, expected string) (bool, error) {
	var isEqual bool
	if val.Type() == cty.String {
		isEqual = cmpString(val, expected)
	} else if val.Type() == cty.Number {
		v, err := strconv.ParseFloat(expected, 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse '%v' to float: %v", expected, err)
		}
		isEqual, err = cmpNumber(val, v)
		if err != nil {
			return false, fmt.Errorf("failed to compare float: %v", err)
		}
	} else {
		return false, fmt.Errorf("cannot handle attributes of type %v", val.Type().FriendlyName())
	}

	return isEqual, nil
}
