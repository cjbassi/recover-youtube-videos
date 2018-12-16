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
    fetchingMissingVideos: false,
    playlists,
  })

export const fetchingMissingVideos = () =>
  action(ActionName.fetchingMissingVideos, { fetchingMissingVideos: true })
