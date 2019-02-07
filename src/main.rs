mod args;
mod youtube;

use structopt::StructOpt;

const CLIENT_SECRET_FILE: &str = "client_secret.json";
const TOKEN_STORE_FILE: &str = "token_store.json";

const LIBRARY_FILE: &str = "library.json";
const RECOVERED_VIDEOS_FILE: &str = "recovered_videos.json";
const UNRECOVERED_VIDEOS_FILE: &str = "unrecovered_videos.json";

const MOCK_API_FILE: &str = "test/mock_api.json";

fn main() {
    let args = args::Args::from_args();

    std::env::set_current_dir(args.directory).unwrap();

    let secret =
        oauth2::read_application_secret(&std::path::PathBuf::from(CLIENT_SECRET_FILE)).unwrap();

    let client = hyper::Client::with_connector(hyper::net::HttpsConnector::new(
        hyper_rustls::TlsClient::new(),
    ));
    let auth = oauth2::Authenticator::new(
        &secret,
        oauth2::DefaultAuthenticatorDelegate,
        client,
        oauth2::DiskTokenStorage::new(&TOKEN_STORE_FILE.to_string()).unwrap(),
        Some(oauth2::FlowType::InstalledInteractive),
    );

    let client = hyper::Client::with_connector(hyper::net::HttpsConnector::new(
        hyper_rustls::TlsClient::new(),
    ));
    let mut hub = youtube3::YouTube::new(client, auth);

    let fetched_library = match args.debug {
        true => json::from_str(&std::fs::read_to_string(MOCK_API_FILE).unwrap()).unwrap(),
        false => youtube::fetch_library(&mut hub).unwrap(),
    };

    let mut non_removed_videos: Vec<youtube::Video> = vec![];
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

    let mut local_library: Vec<youtube::Video> = match std::path::Path::new(LIBRARY_FILE).exists() {
        true => json::from_str(&std::fs::read_to_string(LIBRARY_FILE).unwrap()).unwrap(),
        false => vec![],
    };
    let mut previously_recovered_videos: Vec<youtube::Video> =
        match std::path::Path::new(RECOVERED_VIDEOS_FILE).exists() {
            true => json::from_str::<Vec<youtube::Playlist>>(
                &std::fs::read_to_string(RECOVERED_VIDEOS_FILE).unwrap(),
            )
            .unwrap()
            .into_iter()
            .flat_map(|playlist| {
                playlist
                    .playlist_items
                    .into_iter()
                    .map(|playlist_item| youtube::Video::from(playlist_item))
                    .collect::<Vec<youtube::Video>>()
            })
            .collect(),
            false => vec![],
        };
    local_library.append(&mut previously_recovered_videos);

    let mut playlists_of_recovered_videos: Vec<youtube::Playlist> = vec![];
    let mut playlists_of_unrecovered_videos: Vec<youtube::Playlist> = vec![];
    for playlist in playlists_of_removed_playlist_items {
        let mut playlist_of_recovered_videos: youtube::Playlist = playlist.clone();
        let mut playlist_of_unrecovered_videos: youtube::Playlist = playlist.clone();
        for mut playlist_item in playlist.playlist_items {
            let mut recovered = false;
            for video in &local_library {
                if playlist_item.id == video.id {
                    playlist_item.title = video.title.to_owned();
                    recovered = true;
                    break;
                }
            }
            match recovered {
                true => playlist_of_recovered_videos
                    .playlist_items
                    .push(playlist_item),
                false => playlist_of_unrecovered_videos
                    .playlist_items
                    .push(playlist_item),
            }
        }
        if playlist_of_recovered_videos.playlist_items.len() > 0 {
            playlists_of_recovered_videos.push(playlist_of_recovered_videos)
        }
        if playlist_of_unrecovered_videos.playlist_items.len() > 0 {
            playlists_of_unrecovered_videos.push(playlist_of_unrecovered_videos)
        }
    }

    let j = json::to_string_pretty(&non_removed_videos).unwrap();
    std::fs::write(LIBRARY_FILE, j).unwrap();

    let j = json::to_string_pretty(&playlists_of_recovered_videos).unwrap();
    std::fs::write(RECOVERED_VIDEOS_FILE, j).unwrap();

    let j = json::to_string_pretty(&playlists_of_unrecovered_videos).unwrap();
    std::fs::write(UNRECOVERED_VIDEOS_FILE, j).unwrap();
}
