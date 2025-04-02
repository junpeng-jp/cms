package file

import "errors"

var (
	ErrorMissingInitialFileMarker = errors.New("missing initial file marker")
	ErrorMissingFinalFileMarker   = errors.New("missing final file marker")
	ErrorMalformedFileMetadata    = errors.New("malformed file metadata")
	ErrorMalformedContent         = errors.New("malformed content")

	// Content Decoding Errors

	ErrorInvalidContentRange       = errors.New("invalid content range")
	ErrorInvalidContentStartOffset = errors.New("invalid content start offset")
	ErrorInvalidContentEndOffset   = errors.New("invalid content end offset")
	ErrorContentIsNotValidUTF8     = errors.New("content is not valid utf8")

	// Image Decoding Errors

	ErrorInvalidImageRange       = errors.New("invalid image range")
	ErrorInvalidImageStartOffset = errors.New("invalid image start offset")
	ErrorInvalidImageEndOffset   = errors.New("invalid image end offset")
	ErrorImageIsNotValidUTF8     = errors.New("image is not base64")

	ErrorFileExceedMaxSize                = errors.New("file exceeds max size")
	ErrorUnexpectedOffsetWhenEncodingFile = errors.New("unexpected offset when encoding file")
	ErrorUnableToEncodeFile               = errors.New("unable to encode file")
	ErrorUnexpectedHeadFileMarkerSize     = errors.New("unexpected head marker byte size")
	ErrorUnexpectedTailFileMarkerSize     = errors.New("unexpected tail marker byte size")
)
