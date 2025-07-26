import { TileSet } from "../../TileSet";
import { Rule, State, RuleEvent } from "../../types";

/** Returns a RULE that submits the play areas, scoring them, and moves them to the provided discard pile */
export function applyLetterScores(
  getPlayAreas: (s: State) => TileSet[]
): Rule<"applyLetterScores"> {
  return {
    id: "applyLetterScores",
    label: "rules.applyLetterScores",

    handle: (event: RuleEvent, s: State) => {
      if (event != "submit") {
        return false;
      }

      for (const area of getPlayAreas(s)) {
        for (let i = 1; i < area.length; i++) {
          for (const { handle } of s.rules) {
            const newScore = handle?.(
              "scoreByLetter",
              s,
              area.tiles.slice(0, i)
            );
            if (newScore) {
              s = newScore;
            }
          }
        }
      }
      return s;
    },
  };
}
