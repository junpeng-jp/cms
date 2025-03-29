package encoding

type Transcoder interface {
	Convert([]byte) ([]byte, error)
}

func NewTranscoder(decoder Decoder, encoder Encoder) Transcoder {
	return &transcoder{
		decoder: decoder,
		encoder: encoder,
	}
}

type transcoder struct {
	decoder Decoder
	encoder Encoder
}

func (t *transcoder) Convert(b []byte) ([]byte, error) {
	file, err := t.decoder.Decode(b)
	if err != nil {
		return nil, err
	}

	output, err := t.encoder.Encode(file)
	if err != nil {
		return nil, err
	}

	return output, nil
}
