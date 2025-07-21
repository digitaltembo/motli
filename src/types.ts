import { uuid } from "./uuid";

/** Label used for translation */
export type Label =
  | {
      key: string;
      opts?: Record<string, string>;
    }
  | string;

export type LetterOpts = {
  score?: number;
  specials?: string[];
  empty?: boolean;
};
export class Tile {
  id: string;
  letter: string;
  score: number;
  specials: string[];
  empty: boolean;
  constructor(letter: string, opts?: LetterOpts) {
    this.id = uuid();
    this.letter = letter;
    this.score = opts?.score ?? 0;
    this.specials = opts?.specials ?? [];
    this.empty = opts?.empty ?? false;
  }
  static empty() {
    return new Tile("", { empty: true });
  }

  clone = () => new Tile(this.letter, { ...this });
}

/** Collection of tiles */
export interface TileSet {
  id: string;
  tiles: Tile[];

  /** @returns the total number of tiles in the TileSet (counting duplicates) */
  get length(): number;

  /** Adds the provided letter to the TileSet, returning it */
  add: (l: Tile) => Tile;
  /**
   * Puts a letter at the provided position
   * @param l the letter to add to the TileSet
   * @param position the location within the box to put the letter.
   * @returns the letter that previously was there
   */
  swapIn: (l: Tile, position: number) => Tile;

  /**
   * Swap tiles with the specied other letterbox
   * @param toBox Destination letterbox
   * @param fromLetterId letter id in this letterbox to move to toBox
   * @param toLetterId letter id in toBox to move to this letterbox
   */
  swapWith: (toBox: TileSet, fromLetterId: string, toLetterId: string) => void;
  /**
   * Picks a letter out of the TileSet at random
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the TileSet is empty
   */
  popRandom: (replaceWith?: Tile) => Tile | null;

  /**
   * Picks a letter at a specific index from the TileSet
   * @param letter The specific letter to be removed
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the specified letter isn't in the letter box
   */
  popIndex: (index: number, replaceWith?: Tile) => Tile | null;

  /**
   * Picks a specific letter from the TileSet
   * @param letter The specific letter to be removed
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the specified letter isn't in the letter box
   */
  popSpecific: ({ id }: Tile, replaceWith?: Tile) => Tile | null;

  /** Shuffles the order of the TileSet */
  shuffle: () => void;
  /** Duplicates the TileSet */
  clone: () => TileSet;

  /** Creates a new TileSet of the same size as this one, filled with empty tiles */
  emptyClone: () => TileSet;

  /** Removes all tiles from this TileSet, adding them to the provided TileSet */
  emptyInto: (output: TileSet) => void;
  /** Gets a string representation of the letterbox */
  toString: () => string;

  /**
   * Creates a new TileSet by selecting `size` tiles at random from this letterbox
   * @param size The size of the new TileSet
   * @param replaceWith If provided, will replace the current TileSet with
   *   provided tiles. The passed in letter will be cloned
   * @returns a new letter box
   */
  selectRandom: (size: number, replaceWith?: Tile) => TileSet;
}

/** Store word categories just as numbers, can look up what they mean in the corpus */
export type WordCategory = number;
export type Word = {
  word: string;
  /** Category identifiers */
  categories?: number[];
};

export type Freq = {
  /** Total number of occurrences in the corpus */
  count: number;
  /** Fraction of the corpus */
  fraction: number;
};
/** Analysis of how frequently the given string appears in a corpus overall and at different frequencies within individual words */
export type UsageAnalysis = {
  value: string;
  /** Number of occurrences and percent of the corpus that is this value */
  overallFreq: Freq;
  /**
   * Array of frequencies of multiple occurrences
   * @example for corpus ["ee", "e"] and value "e", this would be
   * [{count: 2, fraction: 1.0}, {count: 1, fraction: 0.5}]
   * meaning "e" occurs at least once in every word,
   * and at least twice in half of the words
   */
  wordFreqs: Freq[];
};

export type Corpus = {
  contains: (s: string) => boolean;
  subCorpus: (predicate: (w: Word) => boolean) => Corpus;
  analyzeUsages: (tiles: Tile[]) => Record<string, UsageAnalysis>;
  categoryLabel: (cat: WordCategory) => Label | null;
};

export class Score {
  basePoints: number;
  multiplier: number;
  bonusPoints: number;

  constructor(s?: Partial<Score>) {
    this.basePoints = s?.basePoints ?? 0;
    this.multiplier = s?.multiplier ?? 0;
    this.bonusPoints = s?.bonusPoints ?? 0;
  }

  get value() {
    return this.basePoints * this.multiplier + this.bonusPoints;
  }

  addToBase = (amount: number) =>
    new Score({ ...this, basePoints: this.basePoints + amount });
  addToBonus = (amount: number) =>
    new Score({ ...this, bonusPoints: this.bonusPoints + amount });
  multiplyMultiplier = (amount: number) =>
    new Score({ ...this, multiplier: this.multiplier * amount });
}

export type RuleEvent = "submit" | "validate" | "scoreLetter" | "scoreWord";
export type Rule<Id extends string = string> = {
  id: Id;
  label: Label;

  scoreByLetter?: (state: State, tiles: Tile[]) => State;
  handle?: (event: RuleEvent, state: State) => State | false;
};

export type Invalidity = {
  label: Label;
  /** If the state is invalid, the provided sections are marked invalid */
  invalidSectionIds?: string[];
};

/** State of the Word Game is stored */
export type State = {
  boxes: TileSet[];
  rules: Rule[];
  corpi: Corpus[];
  /** Time in epoch ms at which the state was last modified */
  lastModifiedTime: number;
  score: Score;
  invalidities: Invalidity[];
  /** Ordered list of partial changes active in this state */
  uncommitted?: State[];
  previous?: State;
};
