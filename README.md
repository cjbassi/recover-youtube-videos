# [recover-youtube-videos](https://recover-youtube-videos.xyz)

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
