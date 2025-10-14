package search

import (
	"clipboard_manager/storage"
	"strings"
)

func LevenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	len1 := len(s1)
	len2 := len(s2)

	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len1][len2]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func FuzzySearch(entries []storage.ClipboardEntry, query string, threshold int) []storage.ClipboardEntry {
	var results []storage.ClipboardEntry

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Text), strings.ToLower(query)) {
			results = append(results, entry)
			continue
		}

		words := strings.Fields(entry.Text)
		for _, word := range words {
			if len(word) < 3 {
				continue
			}
			distance := LevenshteinDistance(query, word)
			if distance <= threshold {
				results = append(results, entry)
				break
			}
		}
	}

	return results
}
