# Corpus Motli

Tools for creating and analyzing word corpora for the purposes of word games.

My goal here is a set of functionality for creating word lists and scoring algorithms for tiles and words that is open-source and easily tunable for different games, play styles, and languages. Scrabble of course has its dictionaries, official and unofficial, in many languages, along with a standardized distribution of tiles and tile scores, but a) they are biased around specific Scrabble rules (e.g. favoring a lot of very short words, not having words longer than 8 or 9 letters) and b) they are _probably_ under copyright, although they do seem to be used widely and are constructed at least partially out of the work of uncompensated volunteers passionate about scrabble.

Written in golang with the idea that it might be a bit more efficient than the TS that is otherwise used throughout Motli at dealing with the large amounts of data.

Built for now using [Wiktionary](https://en.wiktionary.org) as a data source, both as a list of words and as a list of example sentences. May in the future expand to use additional dictionaries (potentially [WordSet](https://github.com/wordset/wordset-dictionary)) and examples of language use (potentially a dump of text from [Wikipedia](https://dumps.wikimedia.org/enwiki/latest/)).

Note that while this code is itself licensed under an [MIT License](./LICENSE), the default dictionaries
used for analysis and therefore the outputted data itself are under the Wiktionary license (CC-BY-SA or GFDL at your choice). The Wiktionary license text can be found at: https://en.wiktionary.org/wiki/Wiktionary:Copyrights.

Outputs data files in the data directory, including fairly large Wiktionary exported files that are downloaded when run.

Available commands:

```

Usage:

	corpus [--download [language]] [--analyze [language [--ngrams [int]] [--tiles [int]]

The flags are:

	--help
			Print the help text

	--download [language]
			Download the wikiextract file for the given language, storing the gzipped jsonl files
			in the data directory

	--analyze [language]
			Run analysis on the language, defaulting to an ngram analysis of size 1

	--analyse [language] --ngrams [int]
			Run ngram analysis of the provided size, storing thee results as a csv in the data directory

	--analyse [language] --tiles [int]
			Run analysis of ngram size of 1 and create a set of tiles of the provided size whose
			frequency corresponds to the frequency of the ngrams in that language's corpus, storing the
			results as a JSON file in the data directory
*/
```

See [NOTES.md](./NOTES.md) for miscellaneous not-fully-categorized thoughts around what went into this, what should go into this, etc

## Related work

- [Alfred Mosher Butts](https://en.wikipedia.org/wiki/Alfred_Mosher_Butts#Scrabble) is the original creator of Scrabble, deduced the frequency of tiles in the tile set based on the frequency that the letters were showing up in the front page of the New York Times (the tiles were originally used for a game which did not have individual scores for each tile). I have not been able to track down a description of how he came up with the scores themselves, other than the fact that he adjusted it over playtesting, and seems obviously related to the frequency of the letters.
- [valette](https://github.com/jmlewis/valett/tree/master) is a project by Joshua Lewis to reassess the value of Scrabble tiles, and the only example I have found of code that actually does it
- [This Deadspin article](https://deadspin.com/h-y-and-z-are-your-money-letters-how-to-beat-scrabbl-5975490/) also reassess the value of Scrabble tiles
