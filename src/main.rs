mod args;
mod utils;
mod youtube;

use std::env;
use std::error::Error;
use std::fs;
use std::path::{Path, PathBuf};
use structopt::StructOpt;

use utils::json_to_file;

const CLIENT_SECRET_FILE: &str = "client_secret.json";
const TOKEN_STORE_FILE: &str = "token_store.json";

const LIBRARY_FILE: &str = "library.json";
const RECOVERED_VIDEOS_FILE: &str = "recovered_videos.json";
const UNRECOVERED_VIDEOS_FILE: &str = "unrecovered_videos.json";

const MOCK_API_FILE: &str = "test/mock_api.json";

type Result<T> = std::result::Result<T, Box<Error>>;

type Hub = google_youtube3::YouTube<
    hyper::Client,
    yup_oauth2::Authenticator<
        yup_oauth2::DefaultAuthenticatorDelegate,
        yup_oauth2::DiskTokenStorage,
        hyper::Client,
    >,
>;

fn create_youtube_hub() -> Result<Hub> {
    let secret = yup_oauth2::read_application_secret(&PathBuf::from(CLIENT_SECRET_FILE))?;

    let client = hyper::Client::with_connector(hyper::net::HttpsConnector::new(
        hyper_rustls::TlsClient::new(),
    ));
    let auth = yup_oauth2::Authenticator::new(
        &secret,
        yup_oauth2::DefaultAuthenticatorDelegate,
        client,
        yup_oauth2::DiskTokenStorage::new(&TOKEN_STORE_FILE.to_string())?,
        Some(yup_oauth2::FlowType::InstalledInteractive),
    );

    let client = hyper::Client::with_connector(hyper::net::HttpsConnector::new(
        hyper_rustls::TlsClient::new(),
    ));
    let hub = google_youtube3::YouTube::new(client, auth);

    Ok(hub)
}

fn get_known_videos() -> Result<Vec<youtube::Video>> {
    let mut local_library: Vec<youtube::Video> = if Path::new(LIBRARY_FILE).exists() {
        serde_json::from_str(&std::fs::read_to_string(LIBRARY_FILE)?)?
    } else {
        vec![]
    };
    let mut previously_recovered_videos: Vec<youtube::Video> =
        if std::path::Path::new(RECOVERED_VIDEOS_FILE).exists() {
            serde_json::from_str::<Vec<youtube::Playlist>>(&fs::read_to_string(
                RECOVERED_VIDEOS_FILE,
            )?)?
            .into_iter()
            .flat_map(|playlist| {
                playlist
                    .playlist_items
                    .into_iter()
                    .map(youtube::Video::from)
                    .collect::<Vec<youtube::Video>>()
            })
            .collect()
        } else {
            vec![]
        };
    local_library.append(&mut previously_recovered_videos);
    Ok(local_library)
}

fn partition_removed_videos(
    fetched_library: Vec<youtube::Playlist>,
) -> (Vec<youtube::Video>, Vec<youtube::Playlist>) {
    let mut non_removed_videos: Vec<youtube::Video> = vec![];
    let mut playlists_of_removed_playlist_items: Vec<youtube::Playlist> = vec![];
    for playlist in fetched_library {
        let mut pl_copy = playlist.clone();
        let (removed, non_removed): (Vec<youtube::PlaylistItem>, Vec<youtube::PlaylistItem>) =
            playlist
                .playlist_items
                .into_iter()
                .partition(|playlist_item| playlist_item.removed());
        non_removed_videos.append(&mut non_removed.into_iter().map(youtube::Video::from).collect());
        if !removed.is_empty() {
            pl_copy.playlist_items = removed;
            playlists_of_removed_playlist_items.push(pl_copy);
        }
    }
    (non_removed_videos, playlists_of_removed_playlist_items)
}

fn partition_recovered_videos(
    playlists_of_removed_playlist_items: Vec<youtube::Playlist>,
    local_library: &[youtube::Video],
) -> (Vec<youtube::Playlist>, Vec<youtube::Playlist>) {
    let mut playlists_of_recovered_videos: Vec<youtube::Playlist> = vec![];
    let mut playlists_of_unrecovered_videos: Vec<youtube::Playlist> = vec![];
    for playlist in playlists_of_removed_playlist_items {
        let mut playlist_of_recovered_videos: youtube::Playlist = playlist.clone();
        let mut playlist_of_unrecovered_videos: youtube::Playlist = playlist.clone();
        for mut playlist_item in playlist.playlist_items {
            let mut recovered = false;
            for video in local_library {
                if playlist_item.id == video.id {
                    playlist_item.title = video.title.to_owned();
                    recovered = true;
                    break;
                }
            }
            if recovered {
                playlist_of_recovered_videos
                    .playlist_items
                    .push(playlist_item)
            } else {
                playlist_of_unrecovered_videos
                    .playlist_items
                    .push(playlist_item)
            }
        }
        if !playlist_of_recovered_videos.playlist_items.is_empty() {
            playlists_of_recovered_videos.push(playlist_of_recovered_videos)
        }
        if !playlist_of_unrecovered_videos.playlist_items.is_empty() {
            playlists_of_unrecovered_videos.push(playlist_of_unrecovered_videos)
        }
    }
    (
        playlists_of_recovered_videos,
        playlists_of_unrecovered_videos,
    )
}

fn main() {
    let args = args::Args::from_args();

    env::set_current_dir(args.directory).unwrap();

    let mut hub = create_youtube_hub().unwrap();

    let fetched_library = if args.debug {
        serde_json::from_str(&fs::read_to_string(MOCK_API_FILE).unwrap()).unwrap()
    } else {
        youtube::fetch_library(&mut hub).unwrap()
    };

    let (non_removed_videos, playlists_of_removed_playlist_items) =
        partition_removed_videos(fetched_library);
    let local_library = get_known_videos().unwrap();
    let (playlists_of_recovered_videos, playlists_of_unrecovered_videos) =
        partition_recovered_videos(playlists_of_removed_playlist_items, &local_library);

    json_to_file(LIBRARY_FILE, &non_removed_videos).unwrap();
    json_to_file(RECOVERED_VIDEOS_FILE, &playlists_of_recovered_videos).unwrap();
    json_to_file(UNRECOVERED_VIDEOS_FILE, &playlists_of_unrecovered_videos).unwrap();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_partition_removed_videos() {}

    #[test]
    fn test_partition_recovered_videos() {}
}
