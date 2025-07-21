/** Returns a RULE that gives a word-level bonus when submitting a word that starts with the given constant */

/** Returns a RULE that gives a word-level bonus when submitting a word that starts with the given constant */
function bonusForEndsWith(endsWith: string, amount: number): Rule {
  return {
    id: `startsWithRule.${endsWith}`,
    label: {
      key: "rules.endsWith",
      opts: { endsWith },
    },
    applyToWord: (tiles, score) => {
      if (tiles.toString().endsWith(endsWith)) {
        score.basePoints += amount;
        return true;
      }
      return false;
    },
  };
}
