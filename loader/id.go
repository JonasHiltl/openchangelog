package loader

type ID int64

func NewID(id int64) ID {
	return ID(id)
}

func (i ID) ToInt() int64 {
	return int64(i)
}
