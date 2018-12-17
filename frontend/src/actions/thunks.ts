import { Dispatch } from 'redux'

import * as backend from '../api/backend'
import * as actions from './actions'
import { Action } from './index'

export const fetchRemovedVideos = (accessToken: any) => {
  return async (dispatch: Dispatch<Action>) => {
    dispatch(actions.fetchingRemovedVideos())
    const playlists = await backend.fetchRemovedVideos(accessToken)
    dispatch(actions.coalescePlaylists(playlists))
  }
}
