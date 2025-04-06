package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/utils"
)

const (
	FileMarker     = "BLOC"
	FileMarkerSize = 4
)

var (
	ErrorProblemDecodingFileMarker = errors.New("problem decoding file marker")
	ErrorMissingInitialFileMarker  = errors.New("missing initial file marker")
	ErrorMissingFinalFileMarker    = errors.New("missing final file marker")
)

func VerifyFileMarkers(reader io.ReadSeeker) error {
	var err error
	b := make([]byte, FileMarkerSize)
	// reading the starting byte marker
	_, err = utils.ReadFromCurrentPosition(reader, b, FileMarkerSize)
	if err != nil {
		return fmt.Errorf("%w:%v", ErrorProblemDecodingFileMarker, err)
	}
	if !bytes.Equal(b, []byte(FileMarker)) {
		// suggests that the file is likely not decodable by this codec
		return fmt.Errorf("%w: got '%s' instead of the expected file marker", ErrorMissingInitialFileMarker, b)
	}
	// reading the ending byte marker
	_, err = utils.ReadFromEndOffset(reader, b, -FileMarkerSize, FileMarkerSize)
	if err != nil {
		return fmt.Errorf("%w:%v", ErrorProblemDecodingFileMarker, err)
	}
	if !bytes.Equal(b, []byte(FileMarker)) {
		// suggests that the file is malformed
		// if the file was read from a stream, then the file may be incomplete
		return fmt.Errorf("%w: got '%s' instead of the expected file marker", ErrorMissingFinalFileMarker, b)
	}
	return nil
}
