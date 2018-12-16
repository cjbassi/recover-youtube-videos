import * as React from 'react'
import { connect } from 'react-redux'

import * as backend from '../api/backend'

interface IVideoList {
  playlists?: backend.IPlaylist[]
}

class VideoList extends React.Component<IVideoList> {
  public render() {
    const { playlists } = this.props
    if (playlists === undefined) {
      return null
    }
    return (
      <ul>
        {playlists.map((playlist: backend.IPlaylist) => (
          <li key={playlist.id}>
            {playlist.title}
            <ul>
              {playlist.videos.map((video: backend.IVideo) => (
                <li key={video.id}>
                  {video.title}
                  <ul>
                    <li key='google url'>
                      <a href={backend.googleURL(video)}>Google</a>
                    </li>
                    <li key='waybackmachine url'>
                      <a href={backend.waybackMachineURL(video)}>Archive.org</a>
                    </li>
                  </ul>
                </li>
              ))}
            </ul>
          </li>
        ))}
      </ul>
    )
  }
}

const mapStateToProps = (state: IVideoList) => {
  return {
    playlists: state.playlists,
  }
}

export default connect(mapStateToProps)(VideoList)
