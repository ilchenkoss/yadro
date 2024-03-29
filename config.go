package main

var (
	siftingTags = map[string]bool{
		"(":    false, //left round bracket
		")":    false, //right round bracket
		",":    false, //comma
		":":    false, //colon
		".":    false, //period
		"''":   false, //closing quotation mark
		"``":   false, //opening quotation mark
		"#":    false, //number sign
		"$":    false, //currency
		"CC":   false, //conjunction, coordinating
		"CD":   false, //cardinal number
		"DT":   false, //determiner
		"EX":   false, //existential there
		"FW":   true,  //foreign word
		"IN":   false, //conjunction, subordinating or preposition
		"JJ":   true,  //adjective
		"JJR":  true,  //adjective, comparative
		"JJS":  true,  //adjective, superlative
		"LS":   true,  //list item marker
		"MD":   false, //verb, modal auxiliary
		"NN":   true,  //noun, singular or mass
		"NNP":  true,  //noun, proper singular
		"NNPS": true,  //noun, proper plural
		"NNS":  true,  //noun, plural
		"PDT":  false, //predeterminer
		"POS":  false, //possessive ending
		"PRP":  false, //pronoun, personal
		"PRP$": false, //pronoun, possessive
		"RB":   true,  //adverb
		"RBR":  true,  //adverb, comparative
		"RBS":  true,  //adverb, superlative
		"RP":   true,  //adverb, particle
		"SYM":  false, //symbol
		"TO":   false, //infinitival to
		"UH":   false, //interjection
		"VB":   true,  //verb, base form
		"VBD":  true,  //verb, past tense
		"VBG":  true,  //verb, gerund or present participle
		"VBN":  true,  //verb, past participle
		"VBP":  true,  //verb, non-3rd person singular present
		"VBZ":  true,  //verb, 3rd person singular present
		"WDT":  false, //wh-determiner
		"WP":   false, //wh-pronoun, personal
		"WP$":  false, //wh-pronoun, possessive
		"WRB":  false, //wh-adverb
	}

	siftingWords = map[string]bool{
		"be":    true,
		"are":   true,
		"am":    true,
		"is":    true,
		"was":   true,
		"were":  true,
		"being": true,
		"can":   true,
		"could": true,
		"do":    true,
		"did":   true,
		"does":  true,
		"doing": true,
		"had":   true,
		"has":   true,
		"may":   true,
		"will":  true,
	}
)
