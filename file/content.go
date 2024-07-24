package file

type (
	Content interface {
		GetBytes() []byte
		Size() int64
	}

	StringContent string
	BytesContent  []byte
)

func (s StringContent) GetBytes() []byte {
	return []byte(s)
}

func (s BytesContent) GetBytes() []byte {
	return s
}

func (s BytesContent) Size() int64 {
	return int64(len(s))
}

func (s StringContent) Size() int64 {
	return int64(len(s))
}
