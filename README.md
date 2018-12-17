# [recover-youtube-videos](https://recover-youtube-videos.xyz)

A webapp that helps you recover privated or deleted videos from your YouTube playlists.

## How it works

1. Read-only permissions are granted by the user
2. All playlists and videos are read
3. Videos are split into 2 groups depending on if they have been removed
4. Non-removed videos are stored in the database for potential later recovery
5. Removed videos are checked against the database for matches and replaced with the match if there is one
6. Removed videos are sent to the client
7. Videos that have been recovered are displayed normally, otherwise links are provided to search Google and Wayback Machine with the URL

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

## TODO

- fetch videos from 'Watch Later'
- testing
- code comments
- client side error handling of backend requests
- fix presentation of videos that have been recovered
- fix router.apply(middlwares)
