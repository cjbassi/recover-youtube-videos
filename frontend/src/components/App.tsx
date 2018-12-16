import * as _ from 'lodash'
import * as React from 'react'
import { connect } from 'react-redux'

import { REPO_URL } from '../constants'
import { IStoreState } from '../store'
import LoggedIn from './LoggedIn'
import LoggedOut from './LoggedOut'

interface IAppProps {
  isLoggedIn: boolean
}

class App extends React.Component<IAppProps> {
  public render() {
    const { isLoggedIn } = this.props
    return (
      <div>
        <h3>
          <a href={REPO_URL}>recover-youtube-videos</a>
        </h3>
        {isLoggedIn ? <LoggedIn /> : <LoggedOut />}
      </div>
    )
  }
}

const mapStateToProps = (state: IStoreState) => {
  return {
    isLoggedIn: state.isLoggedIn,
  }
}

export default connect(mapStateToProps)(App)
