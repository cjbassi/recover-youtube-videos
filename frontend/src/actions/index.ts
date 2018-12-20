import { ActionType } from 'typesafe-actions'

import * as actions from './actions'

export type Action = ActionType<typeof actions>
export * from './actions'
export * from './thunks'

export enum ActionName {
  signedIn = 'SIGNED_IN',
  signedOut = 'SIGNED_OUT',
  fetchedRemovedVideos = 'FETCHED_REMOVED_VIDEOS',
  fetchingRemovedVideos = 'FETCHING_REMOVED_VIDEOS',
  fetchErrored = 'FETCH_ERRORED',
}
