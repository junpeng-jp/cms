package encoding

type Transcoder[I any, O any] interface {
	Convert(I) O
}

type Parser interface {
	Parse([]byte) error
}

type Renderer interface {
	Render() ([]byte, error)
}

type transcoder struct {
	parser Parser
	render Renderer
}
