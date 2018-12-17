import * as _ from 'lodash'
import * as React from 'react'
import { connect } from 'react-redux'
import { Header, Loader } from 'semantic-ui-react'

import { IPlaylist } from '../api/backend'
import { REPO_URL } from '../constants'
import { IStoreState } from '../store'
import LoggedIn from './LoggedIn'
import LoggedOut from './LoggedOut'
import VideoList from './VideoList'

interface IAppProps {
  isLoggedIn: boolean
  fetchingRemovedVideos: boolean
  playlists?: IPlaylist[]
}

class App extends React.Component<IAppProps> {
  public render() {
    const { isLoggedIn, playlists, fetchingRemovedVideos } = this.props
    return (
      <div>
        <div
          className='nav-bar'
          style={{
            display: 'flex',
            justifyContent: 'space-around',
            height: '96px',
            alignItems: 'center',
          }}
        >
          <Header as='h1'>
            <a href={REPO_URL}>recover-youtube-videos</a>
          </Header>
          <div style={{ width: '300px' }}>
            {isLoggedIn ? <LoggedIn /> : <LoggedOut />}
          </div>
        </div>
        {fetchingRemovedVideos && (
          <Loader active={true} inline='centered'>
            Fetching videos
          </Loader>
        )}
        {playlists !== undefined && (
          <div style={{ display: 'flex', justifyContent: 'center' }}>
            <VideoList />
          </div>
        )}
      </div>
    )
  }
}

const mapStateToProps = (state: IStoreState) => {
  return {
    isLoggedIn: state.isLoggedIn,
    fetchingRemovedVideos: state.fetchingRemovedVideos,
    playlists: state.playlists,
  }
}

export default connect(mapStateToProps)(App)
