package jlfuzzy

import (
	"github.com/wneo/goTextDistance"
	"log"
	"sort"
	"strings"
	"unicode/utf8"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

type WordRecord struct {
	word           string       // org words
	mapRuneToCount map[rune]int // rune -> count
	TotolCount     int          // all rune count
}

const (
	AlgorithmLevenshtein        = iota // Levenshtein distance (Default)
	AlgorithmDamerauLevenshtein        // Damerau-Levenshtein distance
	AlgorithmJaro                      // Jaro distance
	AlgorithmJaroWinkler               // Jaro-Winkler distance
)

type JLFuzzy struct {
	mapRunesToCount map[rune]map[int]map[string]*WordRecord
	mapWordToRecord map[string]*WordRecord

	Algorithm int
	EnableLog bool
}

func NewJLFuzzy() *JLFuzzy {
	return &JLFuzzy{
		mapRunesToCount: make(map[rune]map[int]map[string]*WordRecord, 100),
		mapWordToRecord: make(map[string]*WordRecord, 1000),

		Algorithm: AlgorithmLevenshtein,
	}
}

func (j *JLFuzzy) RemoveWords(words []string) {
	for _, word := range words {
		j.RemoveWord(word)
	}
}
func (j *JLFuzzy) RemoveWord(word string) {
	if len(word) == 0 {
		return
	}
	word = strings.ToLower(word)
	if record, ok := j.mapWordToRecord[word]; ok {
		delete(j.mapWordToRecord, word)
		for r, c := range record.mapRuneToCount {
			delete(j.mapRunesToCount[r][c], word)
		}
	}
}

func (j *JLFuzzy) AddWords(words []string) {
	for _, word := range words {
		j.AddWord(word)
	}
}

func (j *JLFuzzy) AddWord(word string) {
	if len(word) == 0 {
		return
	}
	word = strings.ToLower(word)
	if _, ok := j.mapWordToRecord[word]; ok {
		return
	}
	w := j.analysisWorld(word)
	j.mapWordToRecord[word] = w
	for r, c := range w.mapRuneToCount {
		if _, ok := j.mapRunesToCount[r]; !ok {
			j.mapRunesToCount[r] = make(map[int]map[string]*WordRecord, 10)
		}
		if _, ok := j.mapRunesToCount[r][c]; !ok {
			j.mapRunesToCount[r][c] = make(map[string]*WordRecord, 20)
		}

		j.mapRunesToCount[r][c][word] = w
	}

}
func (j *JLFuzzy) analysisWorld(word string) *WordRecord {
	result := make(map[rune]int, len(word))
	totolCount := 0
	for _, v := range word {
		if c, ok := result[v]; ok {
			result[v] = c + 1
		} else {
			result[v] = 1
		}
		totolCount++
	}
	return &WordRecord{word: word, mapRuneToCount: result, TotolCount: totolCount}
}

// word: the word to search
// lack: max count for lack char compare to word.  (0   ~   len(word)-1)
//			eg. 1 for abc  -> ab/bc/ac
//				<0: invalid, auto to be len(word)-1;
//				 0: dont allow lack;
//				>0: allow count;
//
// more: max count for add char compare to word.	(int)
//			eg. 1 for abc  -> abcd
//				<0: no limit;
//				 0: dont allow add;
//				>0: allow count;
// maxCount: max count for results.	(>0)
// minScore: min score for Distance.	(0 ~ 1) 0=all
func (j *JLFuzzy) SearchWord(word string, lack int, more int, maxCount int, minScore float64) (result []string) {
	result = []string{}
	if len(word) == 0 {
		return
	}
	if lack >= len(word) || lack < 0 {
		lack = len(word) - 1
	}
	var orgRecord *WordRecord
	var ok bool
	if orgRecord, ok = j.mapWordToRecord[word]; !ok {
		orgRecord = j.analysisWorld(word)
	}

	// 1. 获取 lack 数以内的 strs

	mapHitLack := make(map[string]int, 20) // 当前命中的缺少字符总数
	allLack := 0                           // 无记录前, 已差的总计数
	for r, c := range orgRecord.mapRuneToCount {

		if allLack > lack && len(mapHitLack) == 0 { // 如果全线已经达到 lack 数, 并无命中纪录, 则直接返回空
			return
		}
		// 1.1 已有 str 全息 + c
		for str, existC := range mapHitLack {
			mapHitLack[str] = existC + c
		}

		// 1.2 遍历
		if mapCountToWords, ok := j.mapRunesToCount[r]; ok {
			for runeCount, mapSTR := range mapCountToWords {
				if runeCount+lack < c { //自身已经不足以满足计数, 直接跳过
					continue
				} else {
					for str, _ := range mapSTR {
						if existC, ok := mapHitLack[str]; ok { // 已有
							if runeCount < c { // 字符数不足
								after := existC - runeCount // 之前已经 +c 过, 这里只需要 - runeCount 即可
								if after > lack {
									delete(mapHitLack, str)
								} else {
									mapHitLack[str] = after
								}
							} else {
								mapHitLack[str] = existC - c // 恢复
							}
						} else { // 新增
							if runeCount < c { // 字符数不足
								after := allLack + c - runeCount
								if after <= lack {
									mapHitLack[str] = after
								}
							} else if allLack <= lack {
								mapHitLack[str] = allLack
							}
						}
					}
				}
			}
		}

		// 1.3 遍历检查, 若有超出计数的, delete
		for str, existC := range mapHitLack {
			if existC > lack {
				delete(mapHitLack, str)
			}
		}

		allLack += c

	}

	// 2. 检查 more
	if more >= 0 {
		for str, lackCount := range mapHitLack {
			record := j.mapWordToRecord[str]
			if record.TotolCount+lackCount-orgRecord.TotolCount > more {
				delete(mapHitLack, str)
			}
		}
	}

	// 3. 预先排除 maxCount, 避免大计算量
	if maxCount > 0 && len(mapHitLack) > maxCount {
		tmpCacheL := make([]int, 0, len(mapHitLack))
		tmpMapCache := make(map[int]string, len(mapHitLack))
		for str, l := range mapHitLack {
			record := j.mapWordToRecord[str]
			v := l + (record.TotolCount + l - orgRecord.TotolCount)
			tmpCacheL = append(tmpCacheL, v)
			tmpMapCache[v] = str
		}
		sort.Sort(sort.IntSlice(tmpCacheL))
		tmpCacheL = tmpCacheL[maxCount:]
		for _, sc := range tmpCacheL {
			if s, ok := tmpMapCache[sc]; ok {
				delete(mapHitLack, s)
			}
		}
	}

	// 4. Jaro–Winkler / Levenshtein distance to order
	scores := make([]float64, 0, len(mapHitLack))
	caches := make(map[float64][]string)
	for str, _ := range mapHitLack {

		var score float64
		switch j.Algorithm {
		case AlgorithmLevenshtein:
			distance := textdistance.LevenshteinDistance(word, str)
			score = 1 - float64(distance)/float64(Max(utf8.RuneCountInString(word), utf8.RuneCountInString(str)))
		case AlgorithmDamerauLevenshtein:
			distance := textdistance.DamerauLevenshteinDistance(word, str)
			score = 1 - float64(distance)/float64(Max(utf8.RuneCountInString(word), utf8.RuneCountInString(str)))
		case AlgorithmJaro:
			score, _ = textdistance.JaroDistance(word, str)
		case AlgorithmJaroWinkler:
			score = textdistance.JaroWinklerDistance(word, str)
		}

		if minScore < 1 && minScore > 0 && score < minScore {
			delete(mapHitLack, str)
			continue
		}
		if _, ok := caches[score]; !ok {
			caches[score] = []string{str}
			scores = append(scores, score)
		} else {
			caches[score] = append(caches[score], str)
		}
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	result = make([]string, 0, len(mapHitLack))
	for _, score := range scores {
		result = append(result, caches[score]...)
	}

	if j.EnableLog {
		log.Println("scores:", scores)
		log.Println("result:", result)
	}

	if maxCount > 0 && len(result) > maxCount {
		result = result[:maxCount]
	}

	return

}

// Max returns the maximum number of passed int slices.
func Max(is ...int) int {
	var max int
	for _, v := range is {
		if max < v {
			max = v
		}
	}
	return max
}
