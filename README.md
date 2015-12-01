# jlfuzzy
fast fuzzy with Levenshtein


search args:

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

usage:

    // 1. create fuzzy module (can save for cache)
    fuzzy := NewJLFuzzy()
    // 2. add words to train
    fuzzy.AddWords([]string{"a", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcd", "bdc", "bcdddd"})
    // 3. search
    result := fuzzy.SearchWord("bdc", 1, -1, 0)
    log.Println(result)
