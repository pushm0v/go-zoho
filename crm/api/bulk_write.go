package api

type ApiBulkWrite interface {
	Write()
}

type apiBulkWrite struct {
	option Option
}

func NewApiBulkWrite(option Option) ApiBulkWrite {
	return &apiBulkWrite{
		option: option,
	}
}

func (bw *apiBulkWrite) Write() {

}
