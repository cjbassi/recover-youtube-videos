mod args;
mod youtube;

use std::env;
use std::error::Error;
use std::fs;
use std::path::{Path, PathBuf};
use structopt::StructOpt;

const CLIENT_SECRET_FILE: &str = "client_secret.json";
const TOKEN_STORE_FILE: &str = "token_store.json";

const LIBRARY_FILE: &str = "library.json";
const RECOVERED_VIDEOS_FILE: &str = "recovered_videos.json";
const UNRECOVERED_VIDEOS_FILE: &str = "unrecovered_videos.json";

const MOCK_API_FILE: &str = "test/mock_api.json";

type BoxResult<T> = Result<T, Box<Error>>;
type Hub = google_youtube3::YouTube<
    hyper::Client,
    yup_oauth2::Authenticator<
        yup_oauth2::DefaultAuthenticatorDelegate,
        yup_oauth2::DiskTokenStorage,
        hyper::Client,
    >,
>;

fn json_to_file<T>(filename: &str, json: &T) -> BoxResult<()>
where
    T: serde::ser::Serialize,
{
    let j = serde_json::to_string_pretty(json)?;
    std::fs::write(filename, j)?;
    Ok(())
}

fn get_youtube_hub() -> BoxResult<Hub> {
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

fn filter_removed_videos(
    fetched_library: youtube::Playlists,
) -> (youtube::Videos, youtube::Playlists) {
    let mut non_removed_videos: youtube::Videos = vec![];
    let mut playlists_of_removed_playlist_items: Vec<youtube::Playlist> = vec![];
    for playlist in fetched_library {
        let mut pl_copy = playlist.clone();
        let (removed, non_removed): (Vec<youtube::PlaylistItem>, Vec<youtube::PlaylistItem>) =
            playlist
                .playlist_items
                .into_iter()
                .partition(|playlist_item| playlist_item.removed());
        non_removed_videos.append(
            &mut non_removed
                .into_iter()
                .map(|playlist_item| youtube::Video::from(playlist_item))
                .collect(),
        );
        if removed.len() > 0 {
            pl_copy.playlist_items = removed;
            playlists_of_removed_playlist_items.push(pl_copy);
        }
    }
    (non_removed_videos, playlists_of_removed_playlist_items)
}

fn get_known_videos() -> BoxResult<youtube::Videos> {
    let mut local_library: youtube::Videos = match Path::new(LIBRARY_FILE).exists() {
        true => serde_json::from_str(&std::fs::read_to_string(LIBRARY_FILE)?)?,
        false => vec![],
    };
    let mut previously_recovered_videos: youtube::Videos =
        match std::path::Path::new(RECOVERED_VIDEOS_FILE).exists() {
            true => serde_json::from_str::<Vec<youtube::Playlist>>(&fs::read_to_string(
                RECOVERED_VIDEOS_FILE,
            )?)?
            .into_iter()
            .flat_map(|playlist| {
                playlist
                    .playlist_items
                    .into_iter()
                    .map(|playlist_item| youtube::Video::from(playlist_item))
                    .collect::<youtube::Videos>()
            })
            .collect(),
            false => vec![],
        };
    local_library.append(&mut previously_recovered_videos);
    Ok(local_library)
}

fn filter_recovered_videos(
    playlists_of_removed_playlist_items: &youtube::Playlists,
    local_library: &youtube::Videos,
) -> (youtube::Playlists, youtube::Playlists) {
    let mut playlists_of_recovered_videos: youtube::Playlists = vec![];
    let mut playlists_of_unrecovered_videos: youtube::Playlists = vec![];
    for playlist in playlists_of_removed_playlist_items {
        let mut playlist_of_recovered_videos: youtube::Playlist = playlist.clone();
        let mut playlist_of_unrecovered_videos: youtube::Playlist = playlist.clone();
        for playlist_item in &playlist.playlist_items {
            let mut clone = (*playlist_item).clone();
            let mut recovered = false;
            for video in local_library {
                if clone.id == video.id {
                    clone.title = video.title.to_owned();
                    recovered = true;
                    break;
                }
            }
            match recovered {
                true => playlist_of_recovered_videos.playlist_items.push(clone),
                false => playlist_of_unrecovered_videos.playlist_items.push(clone),
            }
        }
        if playlist_of_recovered_videos.playlist_items.len() > 0 {
            playlists_of_recovered_videos.push(playlist_of_recovered_videos)
        }
        if playlist_of_unrecovered_videos.playlist_items.len() > 0 {
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

    let mut hub = get_youtube_hub().unwrap();

    let fetched_library = match args.debug {
        true => serde_json::from_str(&fs::read_to_string(MOCK_API_FILE).unwrap()).unwrap(),
        false => youtube::fetch_library(&mut hub).unwrap(),
    };

    let (non_removed_videos, playlists_of_removed_playlist_items) =
        filter_removed_videos(fetched_library);
    let local_library = get_known_videos().unwrap();
    let (playlists_of_recovered_videos, playlists_of_unrecovered_videos) =
        filter_recovered_videos(&playlists_of_removed_playlist_items, &local_library);

    json_to_file(LIBRARY_FILE, &non_removed_videos).unwrap();
    json_to_file(RECOVERED_VIDEOS_FILE, &playlists_of_recovered_videos).unwrap();
    json_to_file(UNRECOVERED_VIDEOS_FILE, &playlists_of_unrecovered_videos).unwrap();
}
