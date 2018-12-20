import { Dispatch } from 'redux'

import * as backend from '../api/backend'
import * as actions from './actions'
import { Action } from './index'

export const fetchRemovedVideos = (accessToken: string) => async (
  dispatch: Dispatch<Action>,
) => {
  dispatch(actions.fetchingRemovedVideos())
  try {
    const removedVideos = await backend.fetchRemovedVideos(accessToken)
    dispatch(actions.fetchedRemovedVideos(removedVideos))
  } catch (e) {
    console.error(e)
    dispatch(actions.fetchErrored())
  }
}
