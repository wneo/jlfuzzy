# jlfuzzy
fast fuzzy with Levenshtein

    // 1. create fuzzy module (can save for cache)
    fuzzy := NewJLFuzzy()
    // 2. add words to train
    fuzzy.AddWords([]string{"a", "abc", "abcd", "aaa", "aaabbb", "ccaa", "bcd", "bdc", "bcdddd"})
    // 3. search
    result := fuzzy.SearchWord("bdc", 1, -1, 0)
    log.Println(result)
