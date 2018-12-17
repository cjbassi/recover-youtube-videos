import { applyMiddleware, createStore } from 'redux'
import thunk from 'redux-thunk'

import { IPlaylist } from './api/backend'
import rootReducer from './reducers'

export interface IStoreState {
  avatarURL?: string
  isLoggedIn: boolean
  userName?: string
  accessToken?: string
  playlists?: IPlaylist[]
  fetchingRemovedVideos: boolean
}

export default createStore<IStoreState, any, any, any>(
  rootReducer,
  applyMiddleware(thunk),
)
