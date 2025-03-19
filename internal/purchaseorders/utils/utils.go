package utils

import "strconv"

// stringOrNil returns a pointer to the string if it's not empty, otherwise nil.
func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntOrNil returns a pointer to the integer if the string is a valid number, otherwise nil.
func IntOrNil(s string) *int {
	if s == "" {
		return nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &value
}

// IntPtrFromInt returns a pointer to the integer.
func IntPtrFromInt(i int) *int {
	return &i
}
