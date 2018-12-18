import { applyMiddleware, createStore } from 'redux'
import thunk from 'redux-thunk'

import { IPlaylist } from './api/backend'
import rootReducer from './reducers'

export interface IStoreState {
  userState?: {
    avatarURL?: string;
    userName?: string;
    accessToken?: string;
  }
  removedVideos?: IPlaylist[]
  fetchingRemovedVideos: boolean
}

export default createStore<IStoreState, any, any, any>(
  rootReducer,
  applyMiddleware(thunk),
)
