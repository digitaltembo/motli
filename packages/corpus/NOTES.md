# Corpus Motli

Miscellaneous thoughts going into this project

A good tile-based word game has 3 core qualities:

- A corpus of canonical words that are valid
- A set of tiles to chose from, with a distribution that lets a representative sample of tiles have enough valid words to be interesting
- Scores assigned to the tiles that correspond in some way to how challenging or interesting it is to play them

## Corpus

Practically probably end up with rules like scrabble

- Proper nouns - probably out, you can find _someone_ named that for a given string of letters, and proper nouns make up a tiny portion of real world use
  - Acronyms are probably a subset of proper nouns
- Foreign language words - probably out, although an n-lingual corpus for people with decent understanding of n-languages might be interesting, ie mix simple words from multiple languages that a particular person/group might understand?
  - Foreign language words that have made it into english parlance seem reasonable
- letters? greek letters? kinda cheaty but 2 letter words are nice?
  - 2 letter words are nice for flexibility when playing in more than 1 dimension. 1 dimensional games probably don't need them?
- common additives? Cutting word construction via appending common prefixes/suffixes is possible, makes the game harder, but would it make it more fun?
  - allow for plurals?
  - allow for -er, -ing, -ion?
- Interjections probably not?

### Reference Corpora:

- [Wiktionary](https://www.wiktionary.org/) - WikiTionary! All of the words!
  - Need to be careful because this is a lot of words, filter out entries with multiple words, entries of foreign language words, entries of proper nouns etc
  - [wikiextract](https://github.com/tatuylonen/wiktextract) is a tool/project to extract and parse Wiktionary into JSONL, with downloads of the results hosted at [kaikki.org](https://kaikki.org/dictionary/rawdata.html), which is apparently a project of a digital archive group, so hopefully has some long-standing endurance
- [WordSet](https://github.com/wordset/wordset-dictionary) - an open source dictionary project, starting from a Princeton project WordNet
  - Not really sure if it has value beyond Wiktionary?
- [SOWPODS](https://en.wikipedia.org/wiki/Collins_Scrabble_Words), also called CSW is the official English list used for Scrabble tournaments outside of Canada and USA
  - Updated every few years, recent https://github.com/scrabblewords/scrabblewords/blob/main/words/British/CSW21.txt
- [NASPA](https://en.wikipedia.org/wiki/NASPA_Word_List), also called NWL, OTCWL, OWL and TWL, is the official English list used for Scrabble tournaments in Canada and USA
  - https://github.com/scrabblewords/scrabblewords/blob/main/words/North-American/NSWL2023.txt
  - Bit of an odd history
    1. The Official Scrabble Players Dictionary (OSPD) removed offensive words in the third edition, called OSPD2(?).
    2. OWL (Official Word List?) was created by adding them back in and adding 9 letter words
    3. OSPD3 added new words to OSPD2, and OWL2? was updated in parallel
    4. Then OTCWL2014 (Official Tournament and Club Word List) was created on top of OWL2, not sure where OSPD went. Then there was a OTCWL2016
    5. The the naming scheme changed again to NWL with NWL2018, NWL2020, and NWL2023
    6. Also TWL was in their somewhere, seems like previous?

## Tile Frequency

There are probably a few different ways of thinking about this. Some qualities we might want out of the distribution of tiles in the bag are:

- the probability we draw a letter _l_ from the bag (called from here-on-out _P(l)_, although to be clear this will be far from formal and I hardly remember my probability and information theory classes) should correspond to the probability that _l_ occurs in a word chosen at random from a uniform distribution of words in our corpus
  - This might be synonymous with _P(l)_ should be defined in such a way to maximize the number of words that can be formed from a given selection of letters taken at random from the distribution. Is it?
- _P(l)_ should correspond to the probability that _l_ occurs in a word chosen at random from words in our corpus _weighted at the frequency that the words in the corpus occur in natural language_
  - Simplest heuristic for this is probably looking at the language comprising the example sentences in the wiktionary. This will have biases not present in natural language, but maybe I can assume they all go away, and just count the occurrences of the letter _l_ in the example sentences in the wiktionary
    - potentially biased towards a uniform distribution of words, as each sentence includes that unique word
    - potentially biased towards phrases common to example-sentence writing
  - Might be worth scanning over all of wikipedia? also has biases?
    - https://dumps.wikimedia.org/enwiki/latest/enwiki-latest-pages-articles-multistream.xml.bz2 19gb compressed data dump
  - Someone has probably precalculated word frequency, maybe ngrams viewer for all words in wiktionary?
  - this might be synonymous with _P(l)_ should be defined in such a way to maximize the number of words an average person will think of given a random selection of letters
- _P(l)_ should be proportional to fun?
  - If there aren't enough es, it will be hard to form words
  * but also, it is fun challenge to form words around the qs, zs, and xs

At the end of the day, it feels like a letter should occur

### Total Tile Count

- It is probably nice to have a clean number of tiles, that make a satisfying sound when dumped on the table at once, that a person can picture in their minds eye and have instincts about what is left/etc
- 'e' is about 88x as common as 'z', so we want a large enough number of tiles to be able to express the distinction in frequencies
- Gut feeling: 75 - 250 is the reasonablish range

### What letters get tiles?

- English is easy - 26 letters, 26 different types of tiles.
  - But What About Digraphs?
    - ng, ch, ck, gh, ph, th, sh, wh, wr, qu?
    - replacing q with qu sounds not terrible
  - and maybe like common digrams?
    - er, re, in
    - he and th is quite common in the corpus, probably cause of the
    - in is more common than hbfyvkwzxq in words
  - or other word groups?
    - ing, tion?
  - a hyphen could be funky
- What about other languages
  - a problem for another day
