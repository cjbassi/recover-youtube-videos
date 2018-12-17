# [recover-youtube-videos](https://recover-youtube-videos.xyz)

A webapp that helps you recover privated or deleted videos from your YouTube playlists.

## Usage

1. Visit the site
2. Login and grant read-only permissions to your YouTube account
3. Click the 'Fetch removed videos` button
4. Peruse the list of videos

All videos are stored in the database on the first run, so future requests will have the database to help with recovering the videos.

## How it works

1. Read-only permissions are granted by the user on the client side
2. The client makes a request to the backend with the granted access_token
3. All playlists and videos are read on the backend
4. Videos are split into 2 groups depending on if they have been removed
5. Non-removed videos are stored in the database for potential later recovery
6. Removed videos are checked against the database for matches and replaced with the match if there is one
7. Removed videos are sent to the client
8. Videos that have been recovered are displayed normally, otherwise links are provided to search Google and Wayback Machine with the URL

## Development

Built with TypeScript, React, Redux, Redux Thunk, Go, and Postgres.  
Backend is deployed with Docker on Heroku and the frontend is hosted with GitHub Pages.

To run locally:

- setup a Google Cloud Platform project with the YouTube API enabled
- copy the API credentials to `backend/client_secrets.json`
- setup the `.env` files with the `CLIENT_ID`
- setup a local Postgres server:
  - `docker run --name ryv-postgres -d -p 5432:5432 postgres`
- start the React development server:
  - `yarn; yarn start`
- start the backend:
  - `go run main.go`
