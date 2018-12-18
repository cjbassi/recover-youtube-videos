import { Dispatch } from 'redux'

import * as backend from '../api/backend'
import * as actions from './actions'
import { Action } from './index'

export const fetchRemovedVideos = (accessToken: string) => async (
  dispatch: Dispatch<Action>,
) => {
  dispatch(actions.fetchingRemovedVideos())
  return backend
    .fetchRemovedVideos(accessToken)
    .then((removedVideos) =>
      dispatch(actions.fetchedRemovedVideos(removedVideos)),
    )
}
