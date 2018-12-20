import * as React from 'react'
import { connect } from 'react-redux'
import { List } from 'semantic-ui-react'

import * as backend from '../api/backend'

interface IVideoList {
  removedVideos: backend.IPlaylist[]
}

class VideoList extends React.Component<IVideoList> {
  public render() {
    const { removedVideos } = this.props
    return (
      <List size={'massive'}>
        {removedVideos.map((playlist: backend.IPlaylist) => (
          <List.Item key={playlist.id}>
            <List.Content>
              <List.Header>{playlist.title}</List.Header>
              <List>
                {playlist.playlistItems.map((video: backend.IPlaylistItem) => (
                  <List.Item key={video.id}>
                    <List.Icon
                      style={{
                        position: 'relative',
                        left: '5px',
                        top: '5px',
                      }}
                      size={'small'}
                    >
                      {video.position}
                    </List.Icon>
                    <List.Content>
                      <div style={{ fontSize: 'large' }}>{video.title}</div>
                      <List horizontal={true}>
                        <List.Item key='google url'>
                          <a href={backend.googleURL(video)}>Google</a>
                        </List.Item>
                        <List.Item key='waybackmachine url'>
                          <a href={backend.waybackMachineURL(video)}>
                            Archive.org
                          </a>
                        </List.Item>
                      </List>
                    </List.Content>
                  </List.Item>
                ))}
              </List>
            </List.Content>
          </List.Item>
        ))}
      </List>
    )
  }
}

const mapStateToProps = (state: IVideoList) => {
  return {
    removedVideos: state.removedVideos,
  }
}

export default connect(mapStateToProps)(VideoList)
