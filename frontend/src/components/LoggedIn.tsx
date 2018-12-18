import * as React from 'react'
import { connect } from 'react-redux'
import { Button, Dropdown, Image } from 'semantic-ui-react'

import * as actions from '../actions'
import { handleAuthClick, revokeAccess } from '../api/google-sign-in'
import { IStoreState } from '../store'

interface ILoggedInProps {
  avatarURL?: string
  userName?: string
  fetchRemovedVideos: () => void
}

class LoggedIn extends React.Component<ILoggedInProps> {
  public render() {
    const { avatarURL, userName, fetchRemovedVideos } = this.props
    return (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-around',
        }}
      >
        <Button onClick={fetchRemovedVideos}>Fetch removed videos</Button>
        <div>
          <Dropdown
            trigger={<Image src={avatarURL} alt='Google avatar image' />}
            icon={null}
            direction='left'
            options={[
              {
                key: 'username',
                text: userName,
                disabled: true,
              },
              {
                key: 'handleAuthClick',
                text: <Button onClick={handleAuthClick}>Log out</Button>,
              },
              {
                key: 'revokeAccess',
                text: <Button onClick={revokeAccess}> Revoke access</Button>,
              },
            ]}
          />
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state: IStoreState) => {
  return {
    ...state.userState,
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
    fetchRemovedVideos: () => {
      dispatch(actions.fetchRemovedVideos(accessToken))
    },
  }
}

export default connect(
  mapStateToProps,
  null,
  mergeProps,
)(LoggedIn)
