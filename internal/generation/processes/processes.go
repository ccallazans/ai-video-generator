package processes

type Process interface {
	Execute(command string) (string, error)
}

type TextProcess interface {
	Process
	generateText(message string) (string, error)
}

type SpeechProcess interface {
	Process
	generateSpeech(message string) (string, error)
}

type VideoProcess interface {
	Process
	generateVideo(message string) (string, error)
}
