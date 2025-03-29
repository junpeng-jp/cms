package storage

import "errors"

var (
	ErrorMissingInitialFileMarker = errors.New("missing initial file marker")
	ErrorMissingFinalFileMarker   = errors.New("missing final file marker")
	ErrorMalformedFileFooter      = errors.New("malfored file footer")
	ErrorMalformedContent         = errors.New("malfored content")

	ErrorFileExceedMaxSize                = errors.New("file exceeds max size")
	ErrorUnexpectedOffsetWhenEncodingFile = errors.New("unexpected offset when encoding file")
	ErrorUnableToEncodeFile               = errors.New("unable to encode file")
	ErrorUnexpectedHeadFileMarkerSize     = errors.New("unexpected head marker byte size")
	ErrorUnexpectedTailFileMarkerSize     = errors.New("unexpected tail marker byte size")
)
