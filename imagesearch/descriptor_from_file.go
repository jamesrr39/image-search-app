package imagesearch

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"image"
	_ "image/gif"  // decode
	_ "image/jpeg" // decode
	_ "image/png"  // decode
	"io"
)

func FileDescriptorFromFileBytes(fileBytes []byte) (*ImageDescriptor, error) {

	picture, _, err := image.Decode(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}

	fileHash, err := HashOfFile(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}
	return NewImageDescriptor(fileHash, picture), nil
}

func HashOfFile(file io.Reader) (string, error) {
	hasher := sha1.New()

	_, err := io.Copy(hasher, file)
	if nil != err {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
