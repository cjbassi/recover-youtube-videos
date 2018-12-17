import * as React from 'react'
import { Button } from 'semantic-ui-react'

import { handleAuthClick } from '../api/google-sign-in'

export default class SignedIn extends React.Component {
  public render() {
    return (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-around',
        }}
      >
        <Button onClick={handleAuthClick}>Log in</Button>
      </div>
    )
  }
}
