import produce from 'immer'

import { Action, ActionName } from './actions'
import { IStoreState } from './store'

const initialState = {
  userState: undefined,
  removedVideos: undefined,
  fetchingRemovedVideos: false,
}

export default function rootReducer(
  state: IStoreState = initialState,
  action: Action,
): IStoreState {
  return produce(state, (draft: IStoreState) => {
    switch (action.type) {
      case ActionName.signedIn:
        draft.userState = action.payload
        return
      case ActionName.signedOut:
        draft.userState = undefined
        return
      case ActionName.fetchingRemovedVideos:
        draft.fetchingRemovedVideos = true
        return
      case ActionName.fetchedRemovedVideos:
        draft.removedVideos = action.payload.removedVideos
        draft.fetchingRemovedVideos = false
        return
    }
  })
}
