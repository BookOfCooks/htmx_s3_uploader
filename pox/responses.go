package pox

type statusOk struct {
	data any
}

type statusNotFound struct{}

func StatusOk(data any) statusOk {
	return statusOk{data: data}
}

func StatusNotFound() statusNotFound {
	return statusNotFound{}
}
