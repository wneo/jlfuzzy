# jlfuzzy
##fast fuzzy with Levenshtein(ignore case, support unicode)

###apis:

#####1. new model -> can save

	NewJLFuzzy()

#####2. add words  -> can call any time
	
	func (j *JLFuzzy) AddWords(words []string) 
	func (j *JLFuzzy) AddWord(word string) 

#####3. remove words  -> can call any time
	
	func (j *JLFuzzy) RemoveWords(words []string)
	func (j *JLFuzzy) RemoveWord(word string)

#####4. config Levenshtein  -> can call any time to update
	
	func (j *JLFuzzy) SetConfig(delCost, insCost, subCost int)


###search args:

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

	// maxScore: max score for Levenshtein.	(>0)

###usage:

    // 1. create fuzzy module (can save for cache)
    fuzzy := NewJLFuzzy()
    // 2. add words to train
    fuzzy.AddWords([]string{"a", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcd", "bdc", "bcdddd"})
    // 3. search
    result := fuzzy.SearchWord("bdc", 1, -1, 0, 100)
    log.Println(result)
