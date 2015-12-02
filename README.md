# jlfuzzy
##fast fuzzy with Levenshtein/Jaro (ignore case, support unicode)

### Implemented
* [Levenshtein distance](http://en.wikipedia.org/wiki/Levenshtein_distance)
* [Damerau-Levenshtein distance](http://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance)
* [Jaro distance](http://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
* [Jaro-Winkler distance](http://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)

### apis:

##### 1. new model -> can save

```go
NewJLFuzzy()
```

##### 2. add words  -> can call any time
	
```go
func (j *JLFuzzy) AddWords(words []string) 
func (j *JLFuzzy) AddWord(word string) 
```

##### 3. remove words  -> can call any time

```go
func (j *JLFuzzy) RemoveWords(words []string)
func (j *JLFuzzy) RemoveWord(word string)
```

##### 4. config Levenshtein  -> can call any time to update

```go
j.Algorithm = ...

const (
	AlgorithmLevenshtein        = iota // Levenshtein distance (Default)
	AlgorithmDamerauLevenshtein        // Damerau-Levenshtein distance
	AlgorithmJaro                      // Jaro distance
	AlgorithmJaroWinkler               // Jaro-Winkler distance
)
```

##### 5. search
```go
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
func (j *JLFuzzy) SearchWord(word string, lack int, more int, maxCount int, minScore float64) (result []string)
```

### How to Use:

```bash
go get github.com/wneo/goTextDistance
go get github.com/wneo/jlfuzzy
```

```go
// 1. create fuzzy module (can save for cache)
fuzzy := NewJLFuzzy()
// 2. add words to train
fuzzy.AddWords([]string{"a", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcd", "bdc", "bcdddd"})
// 3. search
result := fuzzy.SearchWord("bdc", 1, -1, 0, 100)
log.Println(result)
```

### License

This software is released under the MIT License, see LICENSE.txt.
