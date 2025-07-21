import { Rule, RuleEvent, Score, State } from "../../types";
import { addUncommittedState } from "../utils";

/** Score 1 point for each submission */
export const SUBMISSION_SCORE: Rule<"submissionScore"> = {
  id: "submissionScore",
  label: "rules.submissionScore",
  handle: (event: RuleEvent, state: State) => {
    if (event !== "submit") {
      return false;
    }
    return addUncommittedState(state, { score: state.score.addToBase(1) });
  },
};
