# [recover-youtube-videos](https://recover-youtube-videos.xyz)

A webapp that helps you recover privated or deleted videos from your YouTube playlists.

**Note**: doesn't backup 'Watch Later' videos since that resouce isn't provided in the YouTube API

## Usage

1. Visit the site
2. Login and grant read-only permissions to your YouTube account
3. Click the 'Fetch removed videos' button
4. Peruse the list of videos
5. ????
6. Profit

All videos are added to a database on the first run, so future requests will have the database to help with recovering the videos.

## How it works

1. Read-only permissions are granted by the user on the client side
2. The client `POST`s the granted `access_token` to the backend
3. The backend uses the `access_token` to request all playlist videos from the YouTube api
4. The videos are split into 2 groups depending on if they have been removed
5. Non-removed videos are stored in the database for potential later recovery
6. Removed videos are checked against the database for matches and replaced with the match if there is one
7. Removed/recovered videos are returned to the client
8. Videos are presented to the user, and unrecovered videos are supplemented with links to search Google and Wayback Machine for the video

## Development

Built with TypeScript, React, Redux, Redux Thunk, Go, and Postgres.  
Backend is deployed on AWS Lambda using Serverless Framework and the frontend is hosted with GitHub Pages.

## Deployment

1. setup a Google Cloud Platform project with the YouTube API enabled
2. download the API credentials to `backend/client_secrets.json`
3. copy the `client_id` to `frontend/.env.local` and `backend/.env`
4. setup an optional database and copy the URI to `backend/.env`
5. deploy the backend with `make deploy`
6. copy the backend URL to `frontend/.env.local`
7. deploy the frontend with `yarn deploy`
