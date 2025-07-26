import { Tile } from "../../TileSet";
import { Rule, State, RuleEvent } from "../../types";
import { addUncommittedState } from "../utils";

/** Basic rule that adds base points equal to the score of a given letter */
export const LETTER_SCORE: Rule<"letterScore"> = {
  id: "letterScore",
  label: "rules.letterScore",

  handle: (event: RuleEvent, state: State, tiles?: string | Tile[]) => {
    if (event !== "scoreByLetter") {
      return false;
    }
    if (!Array.isArray(tiles)) {
      throw new Error("scoreByLetter event incorrectly thrown");
    }
    return addUncommittedState(state, {
      score: state.score.addToBase(tiles[tiles.length - 1].score),
    });
  },
};
