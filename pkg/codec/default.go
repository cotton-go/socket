package codec

type Default struct {
}

func (Default) Encode(value any) (any, error) {
	return value, nil
}

func (Default) Decode(value any) (any, error) {
	return value, nil
}
