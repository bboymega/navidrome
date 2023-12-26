package scanner2

import (
	"io/fs"
	"time"

	"github.com/navidrome/navidrome/model"
)

type folderEntry struct {
	scanCtx         *scanContext
	path            string    // Full path
	id              string    // DB ID
	updTime         time.Time // From DB
	modTime         time.Time // From FS
	audioFiles      map[string]fs.DirEntry
	imageFiles      map[string]fs.DirEntry
	playlists       []fs.DirEntry
	imagesUpdatedAt time.Time
	tracks          model.MediaFiles
	albums          model.Albums
	artists         model.Artists
	missingTracks   model.MediaFiles
}

func (f *folderEntry) isExpired() bool {
	return f.updTime.Before(f.modTime)
}
