import test, { suite, TestContext } from "node:test";
import { SearchOpts, Tile, TileSet } from "./TileSet";

function testSearch(
  description: string,
  word: string,
  tiles: (Tile | string)[],
  expectedSolutionIdxs: number[][],
  opts?: SearchOpts
) {
  const tileSet = new TileSet(
    tiles.map((t) => (typeof t === "string" ? new Tile(t.split("/")) : t))
  );
  const actual = tileSet.search(word, opts).map((solnTiles) =>
    solnTiles.map(({ id, values }) => ({
      id,
      values,
      idx: tileSet.tiles.findIndex(({ id: toFindId }) => toFindId === id),
    }))
  );
  const expected = expectedSolutionIdxs.map((idxs) =>
    idxs.map((idx) => ({
      id: tileSet.tiles[idx].id,
      values: tileSet.tiles[idx].values,
      idx,
    }))
  );
  test(`Can search ${description}`, (t: TestContext) => {
    // console.log(actual); // easier to see the letters expected
    console.log(actual.map((soln) => soln.map(({ idx }) => idx))); // easier to see the expected solution
    t.assert.strictEqual(actual.length, expected.length);
    t.assert.deepStrictEqual(actual, expected);
  });
}

suite("TileSet Search", () => {
  testSearch("empty", "", [], []);

  testSearch(
    "empty with empty tiles",
    "",
    ["a", Tile.empty(), Tile.empty()],
    []
  );

  testSearch("empty with tiles", "", ["a", "b"], []);
  testSearch("empty with word", "a", [], []);
  testSearch("empty with word and tile", "a", ["z"], []);

  testSearch("single letter", "a", ["a"], [[0]]);

  testSearch("full word", "word", ["w", "o", "r", "d"], [[0, 1, 2, 3]]);
  testSearch(
    "full word with empty tiles",
    "word",
    ["w", "o", "r", "d", Tile.empty(), Tile.empty()],
    [[0, 1, 2, 3]]
  );

  testSearch("repeats", "aaa", ["a"], [[0, 0, 0]]);
  testSearch(
    "multiple solutions",
    "aba",
    ["a", "ba", "b"],
    [
      [0, 1],
      [0, 2, 0],
    ]
  );
  testSearch("unique", "aba", ["a", "ba", "b"], [[0, 1]], {
    requireUnique: true,
  });
  testSearch(
    "combinatorics",
    "aa",
    ["a", "a"],
    [
      [0, 0],
      [1, 0],
      [0, 1],
      [1, 1],
    ]
  );
  testSearch(
    "multiple values on a tile",
    "ab",
    ["a/b", "a/b"],
    [
      [0, 0],
      [1, 0],
      [0, 1],
      [1, 1],
    ]
  );
  // testSearch(
  //   "really long word",
  //   "danceclass",
  //   "ercremchieusatnttoeapcacdcahelhgrmycoosuslreicsa".split(""),
  //   [],
  //   { requireUnique: true } // need to add the neighbor constraint
  // );
});
