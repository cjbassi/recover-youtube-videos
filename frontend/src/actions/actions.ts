import { action } from 'typesafe-actions'

import { IPlaylist } from '../api/backend'
import { ActionName } from './index'

export const signedIn = (
  userName: string,
  avatarURL: string,
  accessToken: string,
) =>
  action(ActionName.signedIn, {
    userName,
    avatarURL,
    accessToken,
  })

export const signedOut = () => action(ActionName.signedOut)

export const fetchErrored = () => action(ActionName.fetchErrored)

export const fetchedRemovedVideos = (removedVideos: IPlaylist[]) =>
  action(ActionName.fetchedRemovedVideos, {
    removedVideos,
  })

export const fetchingRemovedVideos = () =>
  action(ActionName.fetchingRemovedVideos)
