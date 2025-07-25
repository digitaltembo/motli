import { Invalidity, TileSet, State } from "../types";

export function addUncommittedState(
  currentState: State,
  newPartial: Partial<State>
): State {
  const newState = currentState.uncommitted
    ? {
        ...currentState.uncommitted[currentState.uncommitted.length - 1],
        ...newPartial,
      }
    : { ...currentState, ...newPartial };
  return {
    ...currentState,
    uncommitted: currentState.uncommitted
      ? [...currentState.uncommitted, newState]
      : [newState],
  };
}

export function getLetterboxById(
  state: State,
  letterboxId: string
): TileSet | undefined {
  return state.tileSets.find(({ id }) => letterboxId === id);
}

export function addValidationErrors(
  state: State,
  ...invalidities: Invalidity[]
) {
  return {
    ...state,
    invalidities: [...state.invalidities, ...invalidities],
  };
}
