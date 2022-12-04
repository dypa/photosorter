package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"regexp"
	"time"
)

const EXIF_HEADER_READ_SIZE = 12

func ExifParseFileDateTime(filepath string) (time.Time, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return time.Time{}, err
	}
	defer fd.Close()

	return parseExifV4(fd)
}

func parseExifV4(fd *os.File) (time.Time, error) {
	datetime := time.Time{}

	bytesAsString, err := readExifDataAsString(fd)
	if err != nil {
		return time.Time{}, err
	}

	re := regexp.MustCompile(`\d{4}:\d{2}:\d{2}\s\d{2}:\d{2}:\d{2}`)
	if !re.MatchString(bytesAsString) {
		return time.Time{}, err
	} else {
		datetimeStrings := re.FindAllString(bytesAsString, -1)
		for _, v := range datetimeStrings {
			exifDateTime, err := time.ParseInLocation("2006:01:02 15:04:05", v, time.Local)

			if err == nil && !exifDateTime.IsZero() && !exifDateTime.Equal(datetime) && (exifDateTime.Before(datetime) || datetime.IsZero()) {
				datetime = exifDateTime
			}
		}
	}

	return datetime, nil
}

func readExifDataAsString(fd *os.File) (data string, err error) {
	header := make([]byte, 12)
	_, err = fd.Read(header)
	if err != nil {
		return "", err
	}

	if isApp0(header) {
		//trust me i'm engineer
		return readApp1Exif(fd, 32767)
	}

	if isApp1(header) {
		app1DataLength, err := readLength(header)
		if err != nil {
			return "", err
		}

		//-8 means "Please notice that the size "SSSS" includes the size of descriptor itself also."
		return readApp1Exif(fd, app1DataLength-8)
	}

	return "", errors.New("no EXIF in JPEG")
}

func readLength(header []byte) (uint16, error) {
	var length uint16
	err := binary.Read(bytes.NewBuffer(header[4:6]), binary.BigEndian, &length)
	if err != nil {
		return 0, err
	}
	return length, nil
}

// 4 bytes APP0 Marker
// 2 bytes APP0 Data Size
// 5 bytes JFIF Header
func isApp0(header []byte) bool {
	return bytes.Equal(header[0:4], []byte{0xff, 0xd8, 0xff, 0xe0}) && string(header[6:11]) == "JFIF\x00"
}

// 4 bytes APP1 Marker
// 2 bytes APP1 Data Size - SSSS
// 6 bytes Exif Header
func isApp1(header []byte) bool {
	return bytes.Equal(header[0:4], []byte{0xff, 0xd8, 0xff, 0xe1}) && string(header[6:12]) == "Exif\x00\x00"
}

func readApp1Exif(fd *os.File, length uint16) (data string, err error) {
	app1Data := make([]byte, length)
	_, err = fd.Read(app1Data)

	if err != nil {
		return "", err
	}

	return string(app1Data[:]), nil
}
