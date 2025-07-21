import { Tile, Rule, Score, State } from "../../types";
import { addUncommittedState } from "../utils";

/** Basic rule that adds base points equal to the score of a given letter */
export const LETTER_SCORE: Rule<"letterScore"> = {
  id: "letterScore",
  label: "rules.letterScore",

  scoreByLetter: (state: State, tiles: Tile[]) => {
    return addUncommittedState(state, {
      score: state.score.addToBase(tiles[tiles.length - 1].score),
    });
  },
};
