export async function fetchRemovedVideos(
  accessToken: string,
): Promise<IPlaylist[]> {
  if (
    process.env.REACT_APP_BACKEND_API_URL === undefined ||
    process.env.REACT_APP_BACKEND_API_KEY === undefined
  ) {
    throw new Error('.env.production.local file is missing or misconfigured')
  }
  const response = await fetch(process.env.REACT_APP_BACKEND_API_URL, {
    method: 'POST',
    headers: {
      'x-api-key': process.env.REACT_APP_BACKEND_API_KEY,
    },
    mode: 'cors',
    body: JSON.stringify({ access_token: accessToken }),
  })
  return await response.json()
}

export interface IPlaylist {
  id: string
  title: string
  playlistItems: IPlaylistItem[]
}

export interface IPlaylistItem {
  id: string
  title: string
  position: number
}

export function googleURL(playlistItem: IPlaylistItem) {
  return `https://www.google.com/search?q=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3D${
    playlistItem.id
  }`
}

export function waybackMachineURL(playlistItem: IPlaylistItem) {
  return `https://web.archive.org/web/*/https://www.youtube.com/watch?v=${
    playlistItem.id
  }`
}
