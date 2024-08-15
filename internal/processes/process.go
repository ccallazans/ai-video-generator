package processes

type Process interface {
	Execute(request interface{}) (interface{}, error)
	SetNext(handler Process)
}

type GenerationContext struct {
	TempDir        string
	Prompt         string
	Text           string
	SpeechFile     string
	GeneratedVideo string
}
