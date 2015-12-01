package jlfuzzy

import (
	"github.com/wneo/levenshtein/levenshtein"
	"log"
	"sort"
)

var levenshteinOption levenshtein.Options

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	levenshteinOption = levenshtein.DefaultOptions
	levenshteinOption.DelCost = 7
	levenshteinOption.InsCost = 3
	levenshteinOption.SubCost = 5

}

type WordRecord struct {
	word           string       // org words
	mapRuneToCount map[rune]int // rune -> count
	TotolCount     int          // all rune count
}

type JLFuzzy struct {
	mapRunesToCount map[rune]map[int]map[string]*WordRecord
	mapWordToRecord map[string]*WordRecord
}

func NewJLFuzzy() *JLFuzzy {
	return &JLFuzzy{
		mapRunesToCount: make(map[rune]map[int]map[string]*WordRecord, 100),
		mapWordToRecord: make(map[string]*WordRecord, 1000),
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
func (j *JLFuzzy) SearchWord(word string, lack int, more int, maxCount int) (result []string) {
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
			v := l*levenshteinOption.DelCost + (record.TotolCount+l-orgRecord.TotolCount)*levenshteinOption.InsCost
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
	scores := make([]int, 0, len(mapHitLack))
	caches := make(map[int][]string)
	for str, _ := range mapHitLack {
		score := levenshtein.DistanceForStrings([]rune(word), []rune(str), levenshteinOption)
		if _, ok := caches[score]; !ok {
			caches[score] = []string{str}
			scores = append(scores, score)
		} else {
			caches[score] = append(caches[score], str)
		}
	}
	sort.Sort(sort.IntSlice(scores))
	//log.Println("scores:", scores)
	result = make([]string, 0, len(mapHitLack))
	for _, score := range scores {
		result = append(result, caches[score]...)
	}
	if maxCount > 0 && len(result) > maxCount {
		result = result[:maxCount]
	}
	return

}
