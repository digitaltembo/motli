import { Tile, TileSet } from "./types";
import { uuid } from "./uuid";

/** Collection of tiles */
export class BasicTileSet implements TileSet {
  id: string;
  tiles: Tile[];

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

  toString = (): string => {
    let s = "";
    for (const l of this.tiles) {
      s += l.letter;
    }
    return s;
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
}
