// Package convert provides conversion functions.
package convert

import "unsafe"

// B2S converts byte slice to string.
//
//nolint:gosec
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to byte slice.
//
//nolint:gosec
func S2B(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
