import { action } from 'typesafe-actions'

import { IPlaylist } from '../api/backend'
import { ActionName } from '../constants'

export const coalesceGoogleSigninStatus = (
  userName: string,
  avatarURL: string,
  accessToken: string,
) =>
  action(ActionName.coalesceGoogleSigninStatus, {
    isLoggedIn: true,
    userName,
    avatarURL,
    accessToken,
  })

export const coalesceGoogleSignoutStatus = () =>
  action(ActionName.coalesceGoogleSignoutStatus, { isLoggedIn: false })

export const coalescePlaylists = (playlists: IPlaylist[]) =>
  action(ActionName.coalescePlaylists, {
    fetchingRemovedVideos: false,
    playlists,
  })

export const fetchingRemovedVideos = () =>
  action(ActionName.fetchingRemovedVideos, { fetchingRemovedVideos: true })
