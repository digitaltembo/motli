import { Rule, RuleEvent, Score, State } from "../../types";
import { addUncommittedState } from "../utils";

/** Basic rule that adds 1 point for every letter used */
export const LETTER_COUNT_SCORE: Rule<"letterCountScore"> = {
  id: "letterCountScore",
  label: "rules.letterCountScore",

  handle: (event: RuleEvent, state: State) => {
    if (event !== "scoreByLetter") {
      return false;
    }
    return addUncommittedState(state, {
      score: state.score.addToBase(1),
    });
  },
};
