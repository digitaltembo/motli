import { TileSet, Tile } from "../../TileSet";
import { Label, Rule, RuleEvent, State } from "../../types";
import { addValidationErrors } from "../utils";

export function requirePredicate<Id extends string>(
  id: Id,
  label: Label,
  getPlayAreas: (state: State) => TileSet[],
  predicate: (tiles: Tile[]) => boolean
): Rule<`required.${Id}`> {
  return {
    id: `required.${id}`,
    label,

    handle: (event: RuleEvent, state: State) => {
      if (event != "submit") {
        return false;
      }
      for (const area of getPlayAreas(state)) {
        if (!predicate(area.tiles)) {
          return addValidationErrors(state, {
            label: {
              key: `invalid.required.${id}`,
              opts:
                typeof label === "string"
                  ? { labelKey: label }
                  : { ...label.opts, labelKey: label.key },
            },
            invalidSectionIds: [area.id],
          });
        }
      }
      return state;
    },
  };
}
