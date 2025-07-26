import { Corpus, Score, State } from "../types";
import { noRepeats } from "../rules/validators/noRepeats";
import { realWord } from "../rules/validators/realWord";
import { reset } from "../rules/basic/reset";
import { SUBMISSION_SCORE } from "../rules/scores/submissionScore";
import { TileSet } from "../TileSet";

export function boggle(corpus: Corpus, tiles: TileSet): State {
  const playArea = TileSet.empty(16);
  const board = tiles.selectRandom(16);

  const initialState: State = {
    tileSets: [board, playArea],
    corpora: [corpus],
    lastModifiedTime: Date.now(),
    invalidities: [],
    score: new Score(),
    rules: [
      realWord(() => [playArea]),
      noRepeats(() => [playArea]),
      SUBMISSION_SCORE,
      reset(() => ({ [playArea.id]: board.id })),
    ],
  };
  return initialState;
}
