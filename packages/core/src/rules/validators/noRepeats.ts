import { TileSet, Rule, RuleEvent, State } from "../../types";
import { addValidationErrors } from "../utils";

export function noRepeats(
  getPlayAreas: (state: State) => TileSet[],
  allowDuplicatesInSingleSubmission?: boolean
): Rule<"noRepeats"> & { submittedWords: Set<string> } {
  const submittedWords = new Set<string>([]);

  return {
    id: "noRepeats",
    label: {
      key: "rules.noRepeats",
    },
    submittedWords,

    handle: (event: RuleEvent, state: State) => {
      if (event !== "validate") {
        return false;
      }
      const playAreas = getPlayAreas(state);
      for (const area of playAreas) {
        const submittedWord = area.toString();
        if (submittedWords.has(submittedWord)) {
          return addValidationErrors(state, {
            label: {
              key: "invalid.playAreaHasDuplicateWord",
              opts: { duplicateWord: submittedWord },
            },
            invalidSectionIds: [area.id],
          });
        }
        if (!allowDuplicatesInSingleSubmission) {
          submittedWords.add(submittedWord);
        }
      }
      if (allowDuplicatesInSingleSubmission) {
        for (const area of playAreas) {
          submittedWords.add(area.toString());
        }
      }
      return state;
    },
  };
}
