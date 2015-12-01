package jlfuzzy

import (
	"log"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := NewJLFuzzy()
		f.AddWords([]string{"a", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcd", "bdc", "bcdddd"})

		re := f.SearchWord("bdc", 1, -1, 0)
		log.Println(re)
	}
}

var sources = []string{"a", "b", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcdd", "bdc", "bcdddd"}
var testCases = []struct {
	target  string
	lack    int
	more    int
	max     int
	sources []string
	result  []string
}{
	{"", 0xFFFF, -1, 0, sources, []string{}},
	{"abc", 1, -1, 0, sources, []string{"abc", "abcd", "bdc", "bcdd", "aaabbb", "ccaa", "bcdddd"}},
	{"abc", 1, -1, 2, sources, []string{"abc", "abcd"}},
	{"aaa", 0, 0, 0, sources, []string{"aaa"}},
	{"a", 0, 1, 0, sources, []string{"a"}},
	{"aaabbb", 2, -1, 0, sources, []string{"aaabbb"}},
}

func TestSearchString(t *testing.T) {
	for _, testCase := range testCases {
		f := NewJLFuzzy()
		f.AddWords(testCase.sources)
		result := f.SearchWord(testCase.target, testCase.lack, testCase.more, testCase.max)
		//t.Log(result, testCase.result)
		//t.Fail()
		if testEq(result, testCase.result) == false {
			t.Log(
				"Search string",
				testCase.target,
				testCase.lack,
				testCase.more,
				testCase.max,
				"in",
				testCase.sources,
				"failed:",
				result,
				", should be",
				testCase.result)
			t.Fail()
		}
	}
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
