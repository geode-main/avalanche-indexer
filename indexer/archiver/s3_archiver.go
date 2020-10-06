package archiver

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Archiver struct {
	region   string
	bucket   string
	uploader *s3manager.Uploader
}

func NewS3Archiver(region, bucket string) Archiver {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return S3Archiver{
		region:   region,
		bucket:   bucket,
		uploader: s3manager.NewUploader(sess),
	}
}

func (arc S3Archiver) Test() error {
	_, err := arc.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(arc.bucket),
		Key:    aws.String(".avalanche-archiver-test"),
		Body:   strings.NewReader("OK"),
	})
	return err
}

func (arc S3Archiver) Commit(snapshot *Snapshot) error {
	file, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	if err := snapshot.Encode(file); err != nil {
		return err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	key := fmt.Sprintf(
		"%s/%s.json.gz",
		snapshot.Meta.Time.Format("2006-01-02"),
		snapshot.ID,
	)

	tagging := fmt.Sprintf(
		"App=%s&Chain=%s",
		snapshot.Meta.AppName+"/"+snapshot.Meta.AppVersion,
		snapshot.Meta.ChainVersion,
	)

	_, err = arc.uploader.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(arc.bucket),
		Key:     aws.String(key),
		Tagging: aws.String(tagging),
		Body:    file,
	})

	return err
}
