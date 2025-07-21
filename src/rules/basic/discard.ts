import { Rule, State, RuleEvent } from "../../types";
import { getLetterboxById } from "../utils";

/**
 * Function for making a rule that will, on submit, move tiles from "play areas" to "discard piles" on submit
 * @param currentDiscardMap Function returning map of letterbox id of play areas to their respective discard piles
 * @returns
 */
export function discard(
  getDiscardMap: (state: State) => Record<string, string>
): Rule<"discard"> {
  return {
    id: "discard",
    label: "rules.discard",

    handle: (event: RuleEvent, s: State) => {
      if (event != "submit") {
        return false;
      }
      // Clone the current state of the boxes
      const boxes = s.boxes.map((b) => b.clone());

      const discardMap = Object.entries(getDiscardMap(s)).map(
        ([playAreaId, discardId]) => [
          getLetterboxById(s, playAreaId),
          getLetterboxById(s, discardId),
        ]
      );

      for (const [playArea, discard] of discardMap) {
        if (playArea == null || discard == null) {
          throw new Error("Invalid discard pile");
        }
        playArea.emptyInto(discard);
      }

      return {
        ...s,
        lastModifiedTime: Date.now(),
        previous: {
          ...s,
          boxes,
          lastModifiedTime: Date.now(),
        },
      };
    },
  };
}
