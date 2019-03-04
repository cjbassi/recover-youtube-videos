use serde_derive::{Deserialize, Serialize};

use crate::{Hub, Result};

// only 50 items can be returned per youtube query
const MAX_RESULTS: u32 = 50;

#[derive(Serialize, Deserialize)]
pub struct Playlist {
    pub title: String,
    pub id: String,
    pub playlist_items: Vec<PlaylistItem>,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct PlaylistItem {
    pub title: String,
    pub id: String,
    pub position: u32,
}

#[derive(Serialize, Deserialize)]
pub struct Video {
    pub title: String,
    pub id: String,
}

impl Clone for Playlist {
    fn clone(&self) -> Self {
        Playlist {
            title: self.title.to_owned(),
            id: self.id.to_owned(),
            playlist_items: vec![],
        }
    }
}

impl From<google_youtube3::Playlist> for Playlist {
    fn from(playlist: google_youtube3::Playlist) -> Self {
        Playlist {
            title: playlist.snippet.unwrap().title.unwrap(),
            id: playlist.id.unwrap(),
            playlist_items: vec![],
        }
    }
}

impl From<google_youtube3::PlaylistItem> for PlaylistItem {
    fn from(playlist_item: google_youtube3::PlaylistItem) -> Self {
        let snippet = playlist_item.snippet.unwrap();
        PlaylistItem {
            title: snippet.title.unwrap(),
            id: playlist_item.content_details.unwrap().video_id.unwrap(),
            position: snippet.position.unwrap(),
        }
    }
}

impl PlaylistItem {
    pub fn removed(&self) -> bool {
        self.title == "Deleted video" || self.title == "Private video"
    }
}

impl From<PlaylistItem> for Video {
    fn from(playlist_item: PlaylistItem) -> Self {
        Video {
            title: playlist_item.title,
            id: playlist_item.id,
        }
    }
}

fn fetch_playlists(hub: &mut Hub) -> Result<Vec<google_youtube3::Playlist>> {
    let mut page_token = String::new();
    let mut playlists = vec![];
    loop {
        let (_resp, result) = hub
            .playlists()
            .list("snippet")
            .mine(true)
            .max_results(MAX_RESULTS)
            .page_token(&page_token)
            .doit()?;
        playlists.append(&mut result.items.unwrap_or_default());
        match result.next_page_token {
            Some(s) => page_token = s,
            None => break,
        }
    }
    Ok(playlists)
}

fn fetch_playlist_items(
    hub: &mut Hub,
    playlist_id: &str,
) -> Result<Vec<google_youtube3::PlaylistItem>> {
    let mut page_token = String::new();
    let mut playlist_items = vec![];
    loop {
        let (_resp, result) = hub
            .playlist_items()
            .list("snippet,contentDetails")
            .playlist_id(playlist_id)
            .max_results(MAX_RESULTS)
            .page_token(&page_token)
            .doit()?;
        playlist_items.append(&mut result.items.unwrap_or_default());
        match result.next_page_token {
            Some(s) => page_token = s,
            None => break,
        }
    }
    Ok(playlist_items)
}

pub fn fetch_library(hub: &mut Hub) -> Result<Vec<Playlist>> {
    let playlists = fetch_playlists(hub)?
        .into_iter()
        .map(|playlist| {
            let mut new_pl = Playlist::from(playlist);
            new_pl.playlist_items = fetch_playlist_items(hub, &new_pl.id)
                .unwrap()
                .into_iter()
                .map(PlaylistItem::from)
                .collect();
            new_pl
        })
        .collect();
    Ok(playlists)
}
