package platform

// This file contains the implementation of a function for parsing and handling
// string representations of numeric values for RAM and disk space etc

import (
	"github.com/dustin/go-humanize"
)

// parseBytes returns a value for the input string.
//
// This function uses the humanize library from github for go.
//
// Typical inputs can include by way of examples '6gb', '6 GB', '6 GiB'.
// Inputs support SI and IEC sizes.  For more information please review
// https://github.com/dustin/go-humanize/blob/master/bytes.go
//
func ParseBytes(val string) (bytes uint64, err error) {
	return humanize.ParseBytes(val)
}
