package plugin

type YomoObjectPlugin interface {
	Handle(value interface{}) (interface{}, error)
	Observed() string
	Name() string
}

type YomoStreamPlugin interface {
	HandleStream(buf []byte, done bool) ([]byte, error)
	Observed() string
	Name() string
}
