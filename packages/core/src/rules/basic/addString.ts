import { SearchOpts, TileSet, Tile } from "../../TileSet";
import { Rule, RuleEvent, State } from "../../types";
import { addUncommittedState } from "../utils";

type Opts = {
  searchOpts: SearchOpts;
};

export function addString<Id extends string>(
  id: Id,
  getDestination: (s: State) => TileSet,
  getSource: (s: State) => TileSet,
  str: string,
  opts?: Partial<Opts>
): Rule<`addString.${Id}`> {
  return {
    id: `addString.${id}`,
    label: {
      key: "rules.addString",
      opts: { id },
    },
    handle: (event: RuleEvent, s: State, content?: string | Tile[]) => {
      if (event !== "string") {
        return false;
      }
      if (content == null || Array.isArray(content)) {
        throw new Error("Invalid string event thrown with content " + content);
      }

      const destination = getDestination(s).clone();
      const source = getSource(s).clone();
      const tilesToMove = source.search(str, opts?.searchOpts)[0];
      if (tilesToMove.length == 0) {
        throw new Error("Invalid move");
      }

      if (tilesToMove.length > source.tiles.length) {
        throw new Error("Not enough room");
      }
      tilesToMove.map(({ id }, idx) =>
        source.swapWith(destination, id, destination.tiles[idx].id)
      );

      return addUncommittedState(s, {
        tileSets: [
          ...s.tileSets.filter(
            ({ id }) => id !== source.id && id !== destination.id
          ),
          source,
          destination,
        ],
      });
    },
  };
}
