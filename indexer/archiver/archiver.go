package archiver

type Archiver interface {
	Test() error
	Commit(*Snapshot) error
}

var (
	_ Archiver = (*S3Archiver)(nil)
	_ Archiver = (*FileArchiver)(nil)
)
