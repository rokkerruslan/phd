package files

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageListResponse struct {
	EventID  int
	AuthorID int
	Title    string
	Hash     string
	Created  time.Time
}

type ImageUploadRequest struct {
	EventID  int
	AuthorID int
	Title    string
	Data     string

	hash string
}

type ImageRetrieve struct {
	ID       int
	EventID  int
	AuthorID int
	Title    string
	Hash     string
}

func (r *ImageUploadRequest) Validate() error {
	var e []string

	if r.EventID == 0 {
		e = append(e, "`EventID` is missing")
	}
	if r.AuthorID == 0 {
		e = append(e, "`AuthorID` is missing")
	}
	if r.Title == "" {
		e = append(e, "`Title` is missing")
	}
	if r.Data == "" {
		e = append(e, "`Data` is missing")
	}

	if len(e) == 0 {
		return nil
	}

	return fmt.Errorf("ImageUploadRequest.Validate fails: %v", strings.Join(e, ", "))
}

func (r *ImageUploadRequest) Store() error {
	baseErr := "ImageUploadRequest.Store() fails: %v"

	buf, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	hash := sha256.New()
	hash.Write([]byte(r.Title))

	// TODO: more clever solution (use GLOBAL_SALT)
	timeHash, err := time.Now().MarshalBinary()
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	hash.Write(timeHash)
	r.hash = fmt.Sprintf("%x", hash.Sum(nil))

	// TODO: use hash sum for name
	f, err := os.OpenFile(filepath.Join(FilesDirName, r.hash), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	defer f.Close()

	// TODO: fill rewrite file
	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}
