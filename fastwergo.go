package fastwergo

import (
	"fmt"
	"math"
	"strings"
)

// Tokenize splits the input string into tokens based on the char_level and delimiter.
func Tokenize(str string, charLevel bool, delim rune) []string {
	var tokens []string
	if charLevel {
		for _, c := range str {
			tokens = append(tokens, string(c))
		}
	} else {
		tokens = strings.Split(str, string(delim))
	}
	return tokens
}

// RoundToDigits rounds a floating-point number to a specified number of decimal digits.
func RoundToDigits(d float64, digits uint8) float64 {
	if digits >= 7 {
		panic("digits must be less than 7")
	}
	pow10 := []float64{1, 10, 100, 1000, 10000, 100000, 1000000}
	return math.Round(d*pow10[digits]) / pow10[digits]
}

// Compute calculates the edit distance between two strings at either character or word level.
func Compute(hypo, ref string, charLevel bool) (uint64, uint64) {
	hypoTokens := Tokenize(hypo, charLevel, ' ')
	refTokens := Tokenize(ref, charLevel, ' ')

	m := uint64(len(hypoTokens) + 1)
	n := uint64(len(refTokens) + 1)

	f := make([]uint64, m*n)
	for i := uint64(0); i < m; i++ {
		f[i*n] = i
	}
	for j := uint64(0); j < n; j++ {
		f[j] = j
	}

	for i := uint64(1); i < m; i++ {
		for j := uint64(1); j < n; j++ {
			f[i*n+j] = minUint64(f[(i-1)*n+j]+1, f[i*n+(j-1)]+1)
			matchingCase := f[(i-1)*n+(j-1)]
			if hypoTokens[i-1] != refTokens[j-1] {
				matchingCase = matchingCase + uint64(1)
			}
			f[i*n+j] = minUint64(f[i*n+j], matchingCase)
		}
	}

	return f[m*n-1], uint64(len(refTokens))
}

// ScoreSent calculates the error rate for a single pair of hypothesis and reference strings.
func ScoreSent(hypo, ref string, charLevel bool) float64 {
	statsFirst, statsSecond := Compute(hypo, ref, charLevel)
	return RoundToDigits(100*float64(statsFirst)/float64(statsSecond), 4)
}

// Score calculates the error rate for multiple pairs of hypothesis and reference strings.
func Score(hypo, ref []string, charLevel bool) float64 {
	nExamples := len(hypo)
	if nExamples != len(ref) {
		panic("hypo and ref must have the same length")
	}

	var totalEdits, totalLengths float64
	for i := 0; i < nExamples; i++ {
		statsFirst, statsSecond := Compute(hypo[i], ref[i], charLevel)
		totalEdits += float64(statsFirst)
		totalLengths += float64(statsSecond)
	}

	return RoundToDigits(100*totalEdits/totalLengths, 4)
}

// Utility function to find the minimum of two uint64 numbers.
func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	hypo := "the quick brown fox"
	ref := "the quick brown dog"

	// Example usage of ScoreSent
	fmt.Printf("WER: %.4f\n", ScoreSent(hypo, ref, false))

	hypoArr := []string{"the quick brown fox", "hello world"}
	refArr := []string{"the quick brown dog", "hello world"}

	// Example usage of Score
	fmt.Printf("WER: %.4f\n", Score(hypoArr, refArr, false))
}
