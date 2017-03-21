package imagesearch

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"image"
	"image-processing/imageprocessingutil" // todo move this package to shared
	_ "image/gif"                          // decode
	_ "image/jpeg"                         // decode
	_ "image/png"                          // decode
	"io"
	"io/ioutil"

	"github.com/rwcarlsen/goexif/exif"
)

func FileDescriptorFromFile(file io.Reader, qtyBins QtyBins) (*ImageDescriptor, error) {

	fileBytes, err := ioutil.ReadAll(file) // todo scanners
	if nil != err {
		return nil, err
	}

	picture, _, err := image.Decode(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}

	exifData, err := exif.Decode(bytes.NewBuffer(fileBytes))
	if nil == err && nil != exifData {
		pic, err := imageprocessingutil.RotateAndTransformPictureByExifData(picture, *exifData)
		if nil == err {
			picture = pic
		}
	}

	fileHash, err := HashOfFile(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}

	return NewImageDescriptor(fileHash, picture, qtyBins), nil
}

func HashOfFile(file io.Reader) (string, error) {
	hasher := sha1.New()

	_, err := io.Copy(hasher, file)
	if nil != err {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
