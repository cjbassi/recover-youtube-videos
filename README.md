# recover-youtube-videos

A cli app that helps you recover privated and deleted videos from your YouTube playlists.

Creates a json database file on first run that stores the metadata of all playlist videos of a given account, so that any videos that go missing in the future can at least have their name recovered.

**Note**: 'Watch Later' videos are not backed up since that resource isn't provided by the YouTube API.

## Installation

### Prebuilt binaries:

Downloads the correct binary from the releases tab into `$CARGO_HOME/bin`: (currently only x86_64 Linux is prebuilt)

```
bash <(curl https://raw.githubusercontent.com/japaric/trust/c268696ab9f054e1092f195dddeead2420c04261/install.sh) -f --git cjbassi/recover-youtube-videos
```

### From source:

```
cargo install --git https://github.com/cjbassi/recover-youtube-videos
```

## Usage

1. Install the app
2. Setup a Google Cloud Platform project with the YouTube api enabled
3. Create a folder to store the app data
4. Download the API credentials into the folder and name it `client_secret.json`
5. Run the app with the folder path as a cli argument (and authenticate on first run)

Several files will be created:

- `library.json`: acts as a database for playlist videos
- `recovered_videos.json`: videos that have been recovered after checking `library.json`
- `unrecovered_videos.json`: videos that were deleted before `library.json` was created
- `token_store.json`: caches user authorization

**Note**: Unrecovered videos can sometimes be recovered by checking Wayback Machine and Google using the video url.

## Automated run

It's recommended to setup a cron job or systemd timer that runs the program for you, making it so that the library database is kept up to date with newly added videos.

A systemd timer file is located [here](./systemd/recover-youtube-videos.timer) along with the service file [here](./systemd/recover-youtube-videos.service).

To setup the systemd timer:

1. Copy both files to `~/.config/systemd/user/`
2. Replace the command path and cli arg in `recover-youtube-videos.service`
3. Run `systemctl --user daemon-reload`
4. Run `systemctl --user enable --now recover-youtube-videos.timer`
