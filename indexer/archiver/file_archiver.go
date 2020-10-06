package archiver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileArchiver struct {
	dir string
}

func NewFileArchiver(dir string) Archiver {
	return FileArchiver{
		dir,
	}
}

func (arc FileArchiver) Test() error {
	info, err := os.Stat(arc.dir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("archive path is not a directly")
	}

	return nil
}

func (arc FileArchiver) Commit(snapshot *Snapshot) error {
	f, err := os.Create(arc.filename(snapshot))
	if err != nil {
		return err
	}
	defer f.Close()

	return snapshot.Encode(f)
}

func (arc FileArchiver) filename(snapshot *Snapshot) string {
	fullPath := filepath.Join(arc.dir, snapshot.ID)
	return fmt.Sprintf("%s.json.gz", fullPath)
}
