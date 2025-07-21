import { Rule, State, RuleEvent, TileSet } from "../../types";

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
          for (const { scoreByLetter } of s.rules) {
            if (scoreByLetter) {
              s = scoreByLetter(s, area.tiles.slice(0, i));
            }
          }
        }
      }
      return s;
    },
  };
}
