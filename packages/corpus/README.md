# Corpus Motli

Tools for creating and analyzing word corpora for the purposes of word games.

Written in golang with the idea that it might be a bit more efficient than the TS that is otherwise used throughout Motli at dealing with the large amounts of data

Outputs data files in the data directory, including large wiktionary exported files that are downloaded when run

Available commands:

```
./corpus --help # show help
./corpus --download [language, e.g. en] # download wikiextract files to data dir
./corpus --analyze [language, e.g. en] # Analyze letter frequencies within wikiextract example sentences
./corpus --analyze [language, e.g. en] --ngrams [number, e.g. 2] # analyze frequencies of additional ngrams
./corpus --analyze [language, e.g. en] --tiles [number, e.g. 100] # generate distribution of letters into the specified amount of tiles
```

See [NOTES.md](./NOTES.md) for miscellaneous not-fully-categorized thoughts around what went into this, what should go into this, etc
