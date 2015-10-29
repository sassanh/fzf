package algo

import (
	"strings"
	"unicode"

	"github.com/junegunn/fzf/src/util"
)

/*
 * String matching algorithms here do not use strings.ToLower to avoid
 * performance penalty. And they assume pattern runes are given in lowercase
 * letters when caseSensitive is false.
 *
 * In short: They try to do as little work as possible.
 */

func runeAt(runes []rune, index int, max int, forward bool) rune {
	if forward {
		return runes[index]
	}
	return runes[max-index-1]
}

func IsLower(ch rune) bool {
	return ch >= 'a' && ch <= 'z'
}

func IsUpper(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}

func IsLetter(ch rune) bool {
	return IsLower(ch) || IsUpper(ch)
}

func IsDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func IsLetterDigit(ch rune) bool {
	return IsLetter(ch) || IsDigit(ch)
}

func ToLower(char rune) rune {
	if char >= 'A' && char <= 'Z' {
		return char + 32
	}
	return char
}

func FuzzyMatchHelper(runes []rune, r int, pat []rune, p int) (int, int) {
	if p == len(pat) {
		return r, r
	}
	if r == len(runes) {
		return -1, -1
	}

	p_ch := ToLower(pat[p])
	if ToLower(runes[r]) == p_ch {
		_, end := FuzzyMatchHelper(runes, r + 1, pat, p + 1)
		if end >= 0 {
			return p, end
		}
	}
	for i := r + 1; i < len(runes); i++ {
		curr := ToLower(runes[i])
		prev := ToLower(runes[i - 1])
		if curr != p_ch {
			continue
		}
		if (IsUpper(curr) && !IsUpper(prev)) ||
				(IsLower(curr) && !IsLetter(prev)) ||
				(IsDigit(curr) && !IsDigit(prev)) ||
				(!IsLetterDigit(curr) && curr != prev) {
			_, end := FuzzyMatchHelper(runes, i + 1, pat, p + 1)
			if end >= 0 {
				return i, end
			}
		}
	}
	return -1, -1
}

// FuzzyMatch performs fuzzy-match
func FuzzyMatch(caseSensitive bool, forward bool, runes []rune, pattern []rune) (int, int) {
	if len(pattern) == 0 {
		return 0, 0
	}
	return FuzzyMatchHelper(runes, 0, pattern, 0)
}

// ExactMatchNaive is a basic string searching algorithm that handles case
// sensitivity. Although naive, it still performs better than the combination
// of strings.ToLower + strings.Index for typical fzf use cases where input
// strings and patterns are not very long.
//
// We might try to implement better algorithms in the future:
// http://en.wikipedia.org/wiki/String_searching_algorithm
func ExactMatchNaive(caseSensitive bool, forward bool, runes []rune, pattern []rune) (int, int) {
	if len(pattern) == 0 {
		return 0, 0
	}

	lenRunes := len(runes)
	lenPattern := len(pattern)

	if lenRunes < lenPattern {
		return -1, -1
	}

	pidx := 0
	for index := 0; index < lenRunes; index++ {
		char := runeAt(runes, index, lenRunes, forward)
		if !caseSensitive {
			if char >= 'A' && char <= 'Z' {
				char += 32
			} else if char > unicode.MaxASCII {
				char = unicode.To(unicode.LowerCase, char)
			}
		}
		pchar := runeAt(pattern, pidx, lenPattern, forward)
		if pchar == char {
			pidx++
			if pidx == lenPattern {
				if forward {
					return index - lenPattern + 1, index + 1
				}
				return lenRunes - (index + 1), lenRunes - (index - lenPattern + 1)
			}
		} else {
			index -= pidx
			pidx = 0
		}
	}
	return -1, -1
}

// PrefixMatch performs prefix-match
func PrefixMatch(caseSensitive bool, forward bool, runes []rune, pattern []rune) (int, int) {
	if len(runes) < len(pattern) {
		return -1, -1
	}

	for index, r := range pattern {
		char := runes[index]
		if !caseSensitive {
			char = unicode.ToLower(char)
		}
		if char != r {
			return -1, -1
		}
	}
	return 0, len(pattern)
}

// SuffixMatch performs suffix-match
func SuffixMatch(caseSensitive bool, forward bool, input []rune, pattern []rune) (int, int) {
	runes := util.TrimRight(input)
	trimmedLen := len(runes)
	diff := trimmedLen - len(pattern)
	if diff < 0 {
		return -1, -1
	}

	for index, r := range pattern {
		char := runes[index+diff]
		if !caseSensitive {
			char = unicode.ToLower(char)
		}
		if char != r {
			return -1, -1
		}
	}
	return trimmedLen - len(pattern), trimmedLen
}

// EqualMatch performs equal-match
func EqualMatch(caseSensitive bool, forward bool, runes []rune, pattern []rune) (int, int) {
	if len(runes) != len(pattern) {
		return -1, -1
	}
	runesStr := string(runes)
	if !caseSensitive {
		runesStr = strings.ToLower(runesStr)
	}
	if runesStr == string(pattern) {
		return 0, len(pattern)
	}
	return -1, -1
}
