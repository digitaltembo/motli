import { BasicTileSet } from "../BasicTileSet";
import { Corpus, TileSet, Score, State } from "../types";
import { noRepeats } from "../rules/validators/noRepeats";
import { realWord } from "../rules/validators/realWord";
import { reset } from "../rules/basic/reset";
import { SUBMISSION_SCORE } from "../rules/scores/submissionScore";

export function boggle(corpus: Corpus, tiles: TileSet): State {
  const playArea = BasicTileSet.empty(16);
  const board = tiles.selectRandom(16);

  const initialState: State = {
    tileSets: [board, playArea],
    corpi: [corpus],
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
