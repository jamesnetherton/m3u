package m3u

import (
	"fmt"
	"testing"
)

func TestPlaylist(t *testing.T) {
	playlist, _ := Parse("testdata/playlist.m3u")

	if len(playlist.Tracks) != 5 {
		t.Fatalf("Expected track count to be 5")
	}

	for i := 0; i < 5; i++ {
		if playlist.Tracks[i].Length != i+1 {
			t.Fatalf("Expected track Length to be %d but was %d", i+1, playlist.Tracks[i].Length)
		}

		if playlist.Tracks[i].Name != fmt.Sprintf("Track %d", i+1) {
			t.Fatalf("Expected track name to be Track %d but was '%s'", i+1, playlist.Tracks[i].Name)
		}

		if playlist.Tracks[i].URI != fmt.Sprintf("Track%d.mp4", i+1) {
			t.Fatalf("Expected track URI to be Track%d.mp4 but was '%s'", i+1, playlist.Tracks[i].URI)
		}
	}
}

func TestPlaylistInvalidHeader(t *testing.T) {
	_, err := Parse("testdata/playlist_no_header.m3u")
	if err == nil {
		t.Fatalf("Expected parse error")
	}
}

func TestPlaylistBadExtinf(t *testing.T) {
	_, err := Parse("testdata/playlist_bad_extinf_format.m3u")
	if err == nil {
		t.Fatalf("Expected parse error")
	}
}

func TestPlaylistFileNotFound(t *testing.T) {
	_, err := Parse("testdata/playlist_no_exists.m3u")
	if err == nil {
		t.Fatalf("Expected parse error")
	}
}

func TestPlaylistMissingInf(t *testing.T) {
	_, err := Parse("testdata/playlist_missing_inf.m3u")
	if err == nil {
		t.Fatalf("Expected parse error")
	}
}
