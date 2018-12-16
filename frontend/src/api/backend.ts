export async function signIn(idToken: string) {
  await fetch(process.env.REACT_APP_API_URL + '/tokensignin', {
    method: 'POST',
    body: JSON.stringify({ idtoken: idToken }),
    credentials: 'include',
    mode: 'cors',
  })
}

export async function fetchMissingVideos(
  accessToken: string,
): Promise<IPlaylist[]> {
  const response = await fetch(
    process.env.REACT_APP_API_URL + '/fetchmissingvideos',
    {
      method: 'POST',
      body: JSON.stringify({ access_token: accessToken }),
      credentials: 'include',
      mode: 'cors',
    },
  )
  return await response.json()
}

export interface IPlaylist {
  id: string
  title: string
  videos: IVideo[]
}

export interface IVideo {
  id: string
  title: string
}

export function googleURL(video: IVideo) {
  return `https://www.google.com/search?q=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3D${
    video.id
  }`
}

export function waybackMachineURL(video: IVideo) {
  return `https://web.archive.org/web/*/https://www.youtube.com/watch?v=${
    video.id
  }`
}
