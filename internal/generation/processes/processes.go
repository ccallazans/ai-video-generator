package processes

type Process interface {
    Execute(request interface{}) (interface{}, error)
    SetNext(handler Process)
}