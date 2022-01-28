package commander

import (
	"io/fs"
	"os"
	"path/filepath"
)

func hash(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	base := filepath.Base(dir)

	dir = filepath.Dir(dir)

	cab := os.DirFS(dir)

	var infos fileInfos

	err = fs.WalkDir(cab, base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		fi := fileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().Unix(),
		}

		infos = append(infos, fi)

		return nil
	})

	if err != nil {
		return "", err
	}

	return infos.Hash(), nil
}
