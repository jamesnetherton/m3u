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

func TestPlaylistWithStreamInfo(t *testing.T) {
	playlist, _ := Parse("testdata/playlist_stream_info.m3u")

	if len(playlist.VariantStreams) != 4 {
		t.Fatalf("Expected track count to be 4")
	}

	for i := 0; i < 4; i++ {
		if playlist.VariantStreams[i].Bandwidth == 0 {
			t.Fatalf("Expected bandwidth to be 128000 but was %d", playlist.VariantStreams[i].Bandwidth)
		}

		if playlist.VariantStreams[i].Resolution == "" {
			t.Fatalf("Expected resolution to be set but was %s", playlist.VariantStreams[i].Resolution)
		}

		if playlist.VariantStreams[i].AverageBandwith == 0 {
			t.Fatalf("Expected average bandwidth to be set but was %d", playlist.VariantStreams[i].AverageBandwith)
		}

		if playlist.VariantStreams[i].Codecs == "" {
			t.Fatalf("Expected codecs to be set but was %s", playlist.VariantStreams[i].Codecs)
		}

		if playlist.VariantStreams[i].Name == "" {
			t.Fatalf("Expected name to be set but was %s", playlist.VariantStreams[i].Name)
		}

		if playlist.VariantStreams[i].FrameRate == 0 {
			t.Fatalf("Expected frame rate to be set but was %f", playlist.VariantStreams[i].FrameRate)
		}

		if playlist.VariantStreams[i].HdcpLevel == "" {
			t.Fatalf("Expected hdcp level to be set but was %s", playlist.VariantStreams[i].HdcpLevel)
		}

		if playlist.VariantStreams[i].Video == "" {
			t.Fatalf("Expected video to be set but was %s", playlist.VariantStreams[i].Video)
		}

		if playlist.VariantStreams[i].Audio == "" {
			t.Fatalf("Expected audio to be set but was %s", playlist.VariantStreams[i].Audio)
		}

		if playlist.VariantStreams[i].URI == "" {
			t.Fatalf("Expected URI to be set but was %s", playlist.VariantStreams[i].URI)
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
