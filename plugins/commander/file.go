package commander

import (
	"crypto/md5"
	"fmt"
	"io/fs"
)

type CacheFile struct {
	Result Result `json:"result,omitempty"`
	HTML   string `json:"html,omitempty"`
}

type fileInfo struct {
	Name    string      // base name of the file
	Size    int64       // length in bytes for regular files; system-dependent for others
	Mode    fs.FileMode // file mode bits
	ModTime int64       // modification time
	IsDir   bool        // abbreviation for Mode().IsDir()
}

type fileInfos []fileInfo

func (infos fileInfos) Hash() string {
	h := md5.New()
	fmt.Fprint(h, infos)
	hs := fmt.Sprintf("%x", h.Sum(nil))
	return hs
}
