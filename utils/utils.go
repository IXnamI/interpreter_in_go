package utils

import (
	"strconv"
	"strings"
)

func IsLetter(chr byte) bool {
	return ('a' <= chr && chr <= 'z') || ('A' <= chr && chr <= 'Z') || chr == '_'
}

func IsDigit(chr byte) bool {
	return ('0' <= chr) && (chr <= '9')
}

func IsNumberFormat(chr byte) bool {
	return IsDigit(chr) || chr == '.' || chr == 'x' || chr == 'b'
}

func IsFloat(in string) bool {
	if !strings.Contains(in, ".") {
		return false
	}
	_, err := strconv.ParseFloat(in, 64)
	return err == nil
}

func ValidateOctalNotation(in string) bool {
	in = in[2:]
	for _, chr := range in {
		if chr < '0' || chr > '7' {
			return false
		}
	}
	return true
}

func ValidateBinaryNotation(in string) bool {
	in = in[2:]
	for _, chr := range in {
		if chr < '0' || chr > '1' {
			return false
		}
	}
	return true
}
