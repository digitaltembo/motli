import { Tile } from "./TileSet";
import { Word, UsageAnalysis, Corpus, Label, WordCategory } from "./types";
import { uuid } from "./uuid";

/** Count the number of occurrences of value in the provided word */
function occurenceCount(value: string, word: string) {
  let count = 0;
  let pos = 0;
  while (true) {
    pos = word.indexOf(value, pos);
    if (pos >= 0) {
      count++;
      pos++;
    } else {
      break;
    }
  }
  return count;
}

export class BasicCorpus implements Corpus {
  id: string;
  words: Word[];

  constructor(words: Word[]) {
    this.id = uuid();
    this.words = words;
  }

  contains = (s: string) => {
    return this.words.some((w) => w.word === s);
  };

  subCorpus = (predicate: (w: Word) => boolean) => {
    return new BasicCorpus(this.words.filter(predicate));
  };

  categoryLabel = (cat: WordCategory): Label | null => {
    return null;
  };
  analyzeUsages = (tiles: Tile[]): Record<string, UsageAnalysis> => {
    const data: Record<string, UsageAnalysis> = {};
    return data;
    // TODO: should do this probably
    // let maxLength = 0;
    // for(const {word} of this.words) {
    //   if (word.length > maxLength) {
    //     maxLength = word.length;
    //   }
    // }
    // const totalCounts: number[] = Array.from(({length: maxLength}));

    // for (const l of tiles) {
    //   if (l.letter in data) {
    //     data[l.id] = data[l.letter];
    //   } else {
    //     let wordFreqs: Freq[] = Array.from(({length: maxLength}))

    //     for (const word of this.words) {
    //       const count = occurenceCount(l.letter, word.word);
    //       wordFreqs = Array.from({length: count}).map((_, i) => ({
    //         count: (wordFreqs[i]?.count ?? 0) + 1,
    //         frequency: 0
    //       })
    //         if(i < wordFreqs.length) {
    //           return {...wordFre}

    //         }
    //       })
    //       if (count > wordFreqs.length) {

    //       }
    //       data[l.letter] = {
    //         value: l.letter,
    //         overallFreq: {
    //           count,
    //           fraction: 0
    //         },
    //         wordFreqs:
    //       }
    //     }
    //   }
    // }
    // return data;
  };
}
