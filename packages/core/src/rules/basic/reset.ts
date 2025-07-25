import { Rule, State, RuleEvent } from "../../types";
import { getLetterboxById } from "../utils";

/**
 * Function for making a rule that will, on submit, move tiles from "play areas" to "discard piles" on submit
 * @param getResetMap Function returning map of letterbox id of play areas to the boxes that they should be reset to
 * @returns
 */
export function reset(
  getResetMap: (state: State) => Record<string, string>
): Rule<"reset"> {
  return {
    id: "reset",
    label: "rules.reset",

    handle: (event: RuleEvent, s: State) => {
      if (event != "submit") {
        return false;
      }
      // Clone the current state of the boxes
      const tileSets = s.tileSets.map((b) => b.clone());

      const resetMap = Object.entries(getResetMap(s)).map(
        ([playAreaId, discardId]) => [
          getLetterboxById(s, playAreaId),
          getLetterboxById(s, discardId),
        ]
      );

      for (const [from, to] of resetMap) {
        if (from == null || to == null) {
          throw new Error("Invalid discard pile");
        }
        let toLetterIdx = 0;
        for (const l of from.tiles) {
          if (l.empty) {
            if (toLetterIdx > to.tiles.length - 1) {
              throw new Error("Not enough empty tiles to reset");
            }

            while (to.tiles[toLetterIdx]?.empty === false) {
              toLetterIdx += 1;
              if (toLetterIdx > to.tiles.length - 1) {
                throw new Error("Not enough empty tiles to reset");
              }
            }
            from.swapWith(to, l.id, to.tiles[toLetterIdx].id);
            toLetterIdx += 1;
          }
        }
      }

      return {
        ...s,
        lastModifiedTime: Date.now(),
        previous: {
          ...s,
          tileSets,
          lastModifiedTime: Date.now(),
        },
      };
    },
  };
}
