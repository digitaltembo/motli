import { SearchOpts, Tile, TileSet } from "./types";
import { uuid } from "./uuid";

/** I dunno, prevent infinite loops or something */
const ITERATION_CAP = 100_000;
type PartialSolution = [
  /** Tiles used in suffix of this solution */
  soFar: Tile[][],
  /** Unsolved prefix */
  prefix: string,
  /** Solved suffix */
  suffix: string
];
/** Collection of tiles */
export class BasicTileSet implements TileSet {
  id: string;
  tiles: Tile[];

  private solvedSearches: Record<string, Tile[][]> = {};

  constructor(tiles: Tile[] = []) {
    this.id = uuid();
    this.tiles = tiles;
  }

  static empty(length: number) {
    return new BasicTileSet(Array.from({ length }).map(() => Tile.empty()));
  }

  get length() {
    return this.tiles.length;
  }
  add = (l: Tile) => {
    this.tiles.push(l);
    return l;
  };

  /**
   * Puts a letter at the provided position
   * @param l the letter to add to the TileSet
   * @param position the location within the box to put the letter.
   * @returns the letter that previously was there
   */
  swapIn = (l: Tile, position: number): Tile => {
    const currentLetter = this.tiles[position];
    this.tiles[position] = l;
    return currentLetter;
  };

  /**
   * Swap tiles with the specied other letterbox
   * @param toBox Destination letterbox
   * @param fromLetterId letter id in this letterbox to move to toBox
   * @param toLetterId letter id in toBox to move to this letterbox
   */
  swapWith = (toBox: TileSet, fromLetterId: string, toLetterId: string) => {
    const fromLetterIdx = this.tiles.findIndex(({ id }) => id === fromLetterId);
    const toLetterIdx = toBox.tiles.findIndex(({ id }) => id === toLetterId);
    if (fromLetterIdx < 0 || toLetterIdx < 0) {
      throw new Error("Tile did not exist where expected");
    }
    this.swapIn(
      toBox.swapIn(this.tiles[fromLetterIdx], toLetterIdx),
      fromLetterIdx
    );
  };

  /**
   * Picks a letter out of the TileSet at random
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the TileSet is empty
   */
  popRandom = (replaceWith?: Tile): Tile | null => {
    if (this.tiles.length == 0) {
      return null;
    }
    const index = Math.random() * this.tiles.length;
    return this.popSpecific(this.tiles[index], replaceWith);
  };

  /**
   * Picks a letter at a specific index from the TileSet
   * @param letter The specific letter to be removed
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the specified letter isn't in the letter box
   */
  popIndex = (index: number, replaceWith?: Tile): Tile | null => {
    if (index < this.tiles.length) {
      this.popSpecific(this.tiles[index], replaceWith);
    }
    return null;
  };

  /**
   * Picks a specific letter from the TileSet
   * @param letter The specific letter to be removed
   * @param replaceWith If provided, instead of simply removing the letter, it will be replaced with the passed in one
   * @returns the letter that was removed, or null if the specified letter isn't in the letter box
   */
  popSpecific = ({ id }: Tile, replaceWith?: Tile): Tile | null => {
    const specific = this.tiles.find((l) => l.id === id);
    if (specific != null) {
      if (replaceWith) {
        this.tiles = this.tiles.map((l) => (l.id === id ? replaceWith : l));
      } else {
        this.tiles = this.tiles.filter((l) => l.id !== id);
      }
      return specific;
    }
    return null;
  };

  shuffle = () => {
    const newLetters: Tile[] = [];
    let nextLetter: true | Tile | null = true;
    while ((nextLetter = this.popRandom())) {
      newLetters.push(nextLetter);
    }
    this.tiles = newLetters;
  };

  clone = (): TileSet => new BasicTileSet([...this.tiles]);

  /** Creates a new BasicTileSet of the same size as this one, filled with empty tiles */
  emptyClone = (): TileSet =>
    new BasicTileSet(
      Array.from({ length: this.tiles.length }).map(() => Tile.empty())
    );

  /** Removes all tiles from this TileSet, adding them to the provided TileSet */
  emptyInto = (output: TileSet, onlyEmptySlots?: boolean) => {
    if (onlyEmptySlots) {
    }
    this.tiles.forEach((_, i) => output.add(this.swapIn(Tile.empty(), i)));
  };

  /** Hey does there need to be a canonical letter, instead of joining all letters */
  toString = (): string => {
    return this.tiles.flatMap(({ values }) => values).join("");
  };

  selectRandom = (size: number, replaceWith?: Tile): TileSet => {
    const newLetters: Tile[] = [];
    for (let i = 0; i < size; i++) {
      const nextLetter = this.popRandom(replaceWith?.clone());
      if (nextLetter !== null) {
        newLetters.push(nextLetter);
      } else {
        throw new Error("Could not create board");
      }
    }
    return new BasicTileSet(newLetters);
  };

  search = (s: string, opts?: SearchOpts): Tile[][] => {
    const solutions: Record<string, Tile[][]> = { "": [] };

    let iterations = 0;
    let sectionsToSearch = [
      ...new Set([
        ...this.tiles.flatMap(({ values }) => values),
        ...Object.keys(this.solvedSearches),
      ]),
    ];

    const matchesConstraints = (solution: Tile[]): boolean => {
      if (opts?.tileLimit != null && solution.length > opts.tileLimit) {
        return true;
      }
      if (opts?.requireUnique) {
        let idSet = new Set<string>([]);
        for (const { id } of solution) {
          if (idSet.has(id)) {
            return false;
          }
          idSet.add(id);
        }
      }
      if (opts?.customConstraint?.(solution) === false) {
        return false;
      }
      return true;
    };

    while (Object.keys(solutions).length > 0) {
      if (s in this.solvedSearches) {
        return this.solvedSearches[s].filter(matchesConstraints);
      }
      if (s in solutions && opts?.selectFirst) {
        const potentialSolutions = solutions[s].filter(matchesConstraints);
        if (potentialSolutions.length > 0) {
          return potentialSolutions;
        }
      }
      iterations++;
      if (iterations > ITERATION_CAP) {
        // hopefully isn't necessary but like maybe we accidentally do some big searches
        throw new Error("Search failed, too involved");
      }
      // get the shortest prefix
      const [prefix, soFar] = Object.entries(solutions).reduce<
        [string, Tile[][]]
      >(
        (acc, [currentPrefix, currentSolutions]) =>
          currentPrefix.length < acc[0].length
            ? [currentPrefix, currentSolutions]
            : acc,
        [{ length: Number.POSITIVE_INFINITY } as string, []]
      );
      if (prefix) {
        this.solvedSearches[prefix] = soFar;
      }
      delete solutions[prefix];
      const suffix = s.substring(prefix.length);

      // add a solution to the next section on top of the solutions to the previous sections
      const addSolution = (section: string, sectionSolution: Tile[]) => {
        const newPrefix = prefix + section;
        const newSolutions =
          soFar.length === 0
            ? [sectionSolution]
            : soFar.map((singleSolution) => [
                ...singleSolution,
                ...sectionSolution,
              ]);

        if (newPrefix in solutions) {
          solutions[newPrefix].push(...newSolutions);
        } else {
          solutions[newPrefix] = newSolutions;
        }
      };

      for (const section of sectionsToSearch) {
        if (suffix.startsWith(section)) {
          const returnSolutions =
            opts?.selectFirst && section.length === suffix.length;
          if (section in this.solvedSearches) {
            for (const partialSolution of this.solvedSearches[section]) {
              if (returnSolutions) {
                const potentialSolution = [...soFar[0], ...partialSolution];
                if (matchesConstraints(potentialSolution)) {
                  return [potentialSolution];
                }
              }
              addSolution(section, partialSolution);
            }
          } else {
            for (const tile of this.tiles) {
              if (tile.values.includes(section)) {
                if (returnSolutions) {
                  return [[...soFar[0], tile]];
                }
                addSolution(section, [tile]);
              }
            }
          }
        }
      }
    }
    return this.solvedSearches[s]?.filter(matchesConstraints) ?? [];
  };
}
