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
		"FW":   false, //foreign word
		"IN":   false, //conjunction, subordinating or preposition
		"JJ":   true,  //adjective
		"JJR":  false, //adjective, comparative
		"JJS":  false, //adjective, superlative
		"LS":   false, //list item marker
		"MD":   false, //verb, modal auxiliary
		"NN":   true,  //noun, singular or mass
		"NNP":  false, //noun, proper singular
		"NNPS": false, //noun, proper plural
		"NNS":  false, //noun, plural
		"PDT":  false, //predeterminer
		"POS":  false, //possessive ending
		"PRP":  false, //pronoun, personal
		"PRP$": false, //pronoun, possessive
		"RB":   true,  //adverb
		"RBR":  false, //adverb, comparative
		"RBS":  false, //adverb, superlative
		"RP":   false, //adverb, particle
		"SYM":  false, //symbol
		"TO":   false, //infinitival to
		"UH":   false, //interjection
		"VB":   true,  //verb, base form
		"VBD":  false, //verb, past tense
		"VBG":  false, //verb, gerund or present participle
		"VBN":  false, //verb, past participle
		"VBP":  false, //verb, non-3rd person singular present
		"VBZ":  false, //verb, 3rd person singular present
		"WDT":  false, //wh-determiner
		"WP":   false, //wh-pronoun, personal
		"WP$":  false, //wh-pronoun, possessive
		"WRB":  false, //wh-adverb
	}

	siftingWords = map[string]bool{
		"be": true,
	}
)
