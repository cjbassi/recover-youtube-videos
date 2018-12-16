import * as React from 'react'

import { handleAuthClick } from '../api/google-sign-in'

export default class SignedIn extends React.Component {
  public render() {
    return (
      <div>
        <button onClick={handleAuthClick}>Log in</button>
      </div>
    )
  }
}
