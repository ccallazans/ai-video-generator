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
	AspectRatio    string // "16:9", "9:16", "1:1"
	Voice          string // Edge TTS voice name
}
