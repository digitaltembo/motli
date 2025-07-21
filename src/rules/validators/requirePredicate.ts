import { Label, Rule, RuleEvent, State, Tile, TileSet } from "../../types";
import { addUncommittedState } from "../utils";

export function requirePredicate<Id extends string>(
  id: Id,
  label: Label,
  getPlayAreas: (state: State) => TileSet[],
  predicate: (tiles: Tile[]) => boolean
): Rule<`bonus.${Id}`> {
  return {
    id: `bonus.${id}`,
    label,

    handle: (event: RuleEvent, state: State) => {
      if (event != "submit") {
        return false;
      }
      for (const area of getPlayAreas(state)) {
        if (predicate(area.tiles)) {
          state = addUncommittedState(state, {
            score: state.score.addToBase(amount),
          });
        }
      }
      return state;
    },
  };
}
