package m3u

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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

		if playlist.Tracks[i].Tags[0].Name != "group-title" {
			t.Fatalf("Expected tag to be group-title but was '%s'", playlist.Tracks[i].Tags[0].Name)
		}

		if playlist.Tracks[i].Tags[0].Value != "Album1" {
			t.Fatalf("Expected group-title tag value to be Album1 but was '%s'", playlist.Tracks[i].Tags[0].Value)
		}
	}
}

func TestRemotePlaylist(t *testing.T) {
	playlist, _ := Parse("https://raw.githubusercontent.com/jamesnetherton/m3u/master/testdata/playlist.m3u")

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

		if playlist.Tracks[i].Tags[0].Name != "group-title" {
			t.Fatalf("Expected tag to be group-title but was '%s'", playlist.Tracks[i].Tags[0].Name)
		}

		if playlist.Tracks[i].Tags[0].Value != "Album1" {
			t.Fatalf("Expected group-title tag value to be Album1 but was '%s'", playlist.Tracks[i].Tags[0].Value)
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

func TestMarshallPlaylist(t *testing.T) {
	playlist, err := Parse("testdata/playlist.m3u")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := Marshall(playlist)
	if err != nil {
		t.Fatal(err)
	}

	b := reader.(*bytes.Buffer)

	ioutil.WriteFile("/tmp/test_m3u_marshalling.m3u", b.Bytes(), os.ModePerm)

	playlist, _ = Parse("/tmp/test_m3u_marshalling.m3u")

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

		if playlist.Tracks[i].Tags[0].Name != "group-title" {
			t.Fatalf("Expected tag to be group-title but was '%s'", playlist.Tracks[i].Tags[0].Name)
		}

		if playlist.Tracks[i].Tags[0].Value != "Album1" {
			t.Fatalf("Expected group-title tag value to be Album1 but was '%s'", playlist.Tracks[i].Tags[0].Value)
		}
	}
}
