import { Action } from './actions'
import { ActionName } from './constants'
import { IStoreState } from './store'

const initialState = {
  avatarURL: undefined,
  isLoggedIn: false,
  userName: undefined,
  accessToken: undefined,
  playlists: undefined,
  fetchingMissingVideos: false,
}

export default function rootReducer(
  state: IStoreState = initialState,
  action: Action,
): IStoreState {
  switch (action.type) {
    case ActionName.coalesceGoogleSigninStatus:
      return {
        ...state,
        ...action.payload,
      }
    case ActionName.coalesceGoogleSignoutStatus:
      return {
        ...state,
        ...action.payload,
      }
    case ActionName.coalescePlaylists:
      return {
        ...state,
        ...action.payload,
      }
    case ActionName.fetchingMissingVideos:
      return {
        ...state,
        ...action.payload,
      }
    default:
      return state
  }
}
