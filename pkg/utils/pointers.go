package utils

// IntPtr returns a pointer to an int literal
func IntPtr(v int) *int {
	return &v
}

// FloatPtr returns a pointer to a float64 literal
func FloatPtr(v float64) *float64 {
	return &v
}
