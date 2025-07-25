import { Invalidity, TileSet, Rule, RuleEvent, State } from "../../types";
import { addValidationErrors } from "../utils";

function validateState(
  currentState: State,
  playAreas: TileSet[]
): Invalidity | undefined {
  const invalidPlayAreas = playAreas.flatMap((playArea) => {
    const word = playArea.toString();
    const valid = currentState.corpi.some((c) => c.contains(word));
    return valid ? [] : [playArea.id];
  });
  if (invalidPlayAreas.length) {
    return {
      label: "invalid.playAreaIsNotRealWord",
      invalidSectionIds: invalidPlayAreas,
    };
  }
  return undefined;
}

export function realWord(
  getPlayAreas: (state: State) => TileSet[]
): Rule<"realWord"> {
  return {
    id: "realWord",
    label: "rules.realWord",

    handle: (event: RuleEvent, currentState: State) => {
      if (event !== "validate") {
        return false;
      }
      const invalidWord = validateState(
        currentState,
        getPlayAreas(currentState)
      );
      if (invalidWord) {
        return addValidationErrors(currentState, invalidWord);
      }
      return currentState;
    },
  };
}
