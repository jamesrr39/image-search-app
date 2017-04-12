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
	"io/ioutil"

	"github.com/jamesrr39/goutil/image-processing/imageprocessingutil"

	"github.com/rwcarlsen/goexif/exif"
)

func FileDescriptorFromFile(file io.Reader, qtyBins QtyBins, location Location, maxGoRoutines uint) (*PersistedImageDescriptor, error) {

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

	return NewPersistedImageDescriptor(NewImageDescriptorFromPicture(fileHash, picture, qtyBins, maxGoRoutines), location), nil
}

func HashOfFile(file io.Reader) (string, error) {
	hasher := sha1.New()

	_, err := io.Copy(hasher, file)
	if nil != err {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
