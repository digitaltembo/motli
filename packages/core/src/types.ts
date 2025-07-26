import { Tile, TileSet } from "./TileSet";

/** Label used for translation */
export type Label =
  | {
      key: string;
      opts?: Record<string, string>;
    }
  | string;

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

export type RuleEvent =
  /** Thrown to validate prior to submission */
  | "validate"
  /** Thrown when the play areas are submitted */
  | "submit"
  /**
   * Thrown to calculate the score, one letter at a time.
   * The content will be the full selection to enable greater contextual knowledge,
   * but the score should only reflect the marginal change to the score incurred by the
   * last character in the list
   **/
  | "scoreByLetter"
  /** Thrown when a single character was selected */
  | "character"
  /** Thrown by game when a string of characters is selected */
  | "string";
export type Rule<Id extends string = string> = {
  id: Id;
  label: Label;

  handle?: (
    event: RuleEvent,
    state: State,
    content?: string | Tile[]
  ) => State | false;
};

export type Invalidity = {
  label: Label;
  /** If the state is invalid, the provided sections are marked invalid */
  invalidSectionIds?: string[];
};

/** State of the Word Game is stored */
export type State = {
  tileSets: TileSet[];
  rules: Rule[];
  corpora: Corpus[];
  /** Time in epoch ms at which the state was last modified */
  lastModifiedTime: number;
  score: Score;
  invalidities: Invalidity[];
  /** Ordered list of partial changes active in this state */
  uncommitted?: State[];
  previous?: State;
};
