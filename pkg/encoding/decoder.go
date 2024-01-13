package encoding

type Decoder interface {
	Decode(any) error
}

type Encoder interface {
	Encode(any) error
}
