import * as React from 'react'
import { connect } from 'react-redux'

import * as actions from '../actions'
import { handleAuthClick, revokeAccess } from '../api/google-sign-in'
import { IStoreState } from '../store'
import VideoList from './VideoList'

interface ILoggedInProps {
  avatarURL?: string
  userName?: string
  fetchingMissingVideos: boolean
  fetchMissingVideos: () => void
}

class LoggedIn extends React.Component<ILoggedInProps> {
  public render() {
    const {
      avatarURL,
      userName,
      fetchMissingVideos,
      fetchingMissingVideos,
    } = this.props
    return (
      <div>
        <button onClick={revokeAccess}>Revoke access</button>
        <button onClick={handleAuthClick}>Log out</button>
        <button onClick={fetchMissingVideos}>Fetch missing videos</button>
        <div>
          <div>{userName}</div>
          <img src={avatarURL} alt='Google avatar image' />
        </div>
        {fetchingMissingVideos && <div>Fetching playlists</div>}
        <VideoList />
      </div>
    )
  }
}

const mapStateToProps = (state: IStoreState) => {
  return {
    avatarURL: state.avatarURL,
    userName: state.userName,
    accessToken: state.accessToken,
    fetchingMissingVideos: state.fetchingMissingVideos,
  }
}

const mergeProps = (
  stateProps: any,
  dispatchProps: any,
  ownProps: any,
): ILoggedInProps => {
  const { dispatch } = dispatchProps
  const { accessToken } = stateProps
  return {
    ...stateProps,
    ...ownProps,
    fetchMissingVideos: () => {
      dispatch(actions.fetchMissingVideos(accessToken))
    },
  }
}

export default connect(
  mapStateToProps,
  null,
  mergeProps,
)(LoggedIn)
