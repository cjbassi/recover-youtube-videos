import * as actions from '../actions'
import Store from '../store'

const YOUTUBE_SCOPE = 'https://www.googleapis.com/auth/youtube.readonly'
const DISCOVERY_URLS = [
  'https://www.googleapis.com/discovery/v1/apis/youtube/v3/rest',
]

export let gapi: any
let GoogleAuth: any

// tslint:disable-next-line:no-var-requires
require('google-client-api')().then((gApi: any) => {
  gapi = gApi
  gapi.load('client:auth2', initClient)
})

async function initClient() {
  await gapi.client.init({
    clientId: process.env.REACT_APP_CLIENT_ID,
    discoveryDocs: DISCOVERY_URLS,
    scope: YOUTUBE_SCOPE,
  })
  GoogleAuth = gapi.auth2.getAuthInstance()
  GoogleAuth.isSignedIn.listen(dispatchSigninStatusUpdate)
  dispatchSigninStatusUpdate()
}

function dispatchSigninStatusUpdate() {
  const googleUser = GoogleAuth.currentUser.get()
  const isAuthorized = googleUser.hasGrantedScopes(YOUTUBE_SCOPE)
  if (isAuthorized) {
    const authResponse = googleUser.getAuthResponse()
    const accessToken = authResponse.access_token
    const idToken = authResponse.id_token
    Store.dispatch(actions.logginBackend(idToken))

    const userProfile = googleUser.getBasicProfile()
    const userName = userProfile.getName()
    const imageUrl = userProfile.getImageUrl()
    Store.dispatch(
      actions.coalesceGoogleSigninStatus(userName, imageUrl, accessToken),
    )
  } else {
    Store.dispatch(actions.coalesceGoogleSignoutStatus())
  }
}

export function handleAuthClick() {
  if (GoogleAuth.isSignedIn.get()) {
    GoogleAuth.signOut()
  } else {
    GoogleAuth.signIn()
  }
}

export function revokeAccess() {
  GoogleAuth.disconnect()
}
