import { Dispatch } from 'redux'

import * as backend from '../api/backend'
import * as actions from './actions'
import { Action } from './index'

export const logginBackend = (idToken: any) => {
  return async (dispatch: Dispatch<Action>) => {
    await backend.signIn(idToken)
  }
}

export const fetchMissingVideos = (accessToken: any) => {
  return async (dispatch: Dispatch<Action>) => {
    dispatch(actions.fetchingMissingVideos())
    const playlists = await backend.fetchMissingVideos(accessToken)
    dispatch(actions.coalescePlaylists(playlists))
  }
}
