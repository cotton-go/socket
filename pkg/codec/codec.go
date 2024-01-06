package codec

type ICodec interface {
	Encode(value any) (any, error)
	Decode(value any) (any, error)
}
