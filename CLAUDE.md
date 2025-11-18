# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI Video Story Generator that transforms text prompts into narrated videos with auto-generated captions. The system uses Ollama LLM for text generation, AI text-to-speech for narration, and Whisper for caption generation.

## Technology Stack

- **Backend**: Go 1.21.1 with Echo web framework
- **LLM**: Ollama (default model: orca-mini)
- **TTS**: Facebook MMS-TTS (transformers library)
- **Caption Generation**: OpenAI Whisper
- **Video Processing**: FFmpeg, MoviePy
- **Deployment**: Docker Compose with multi-container setup

## Running the Project

### Docker (Recommended)
```bash
make start                    # Start all services via docker-compose
```

First run downloads ~5.5GB of Ollama models and Python dependencies. Subsequent runs use Docker layer caching.

### Local Development
```bash
# 1. Install and run Ollama
# Download from: https://ollama.com/download

# 2. Setup Python environment
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# 3. Run Go application
go run cmd/*
```

### Testing the API
```bash
curl --location 'http://localhost:8080/api/v1/generate' \
--header 'Content-Type: application/json' \
--data '{"message": "Tell me a story about Bahia"}'
```

Generated videos are saved in the `./generated` folder.

## Architecture

### High-Level Flow
The application implements a **Chain of Responsibility** pattern through the `Process` interface (internal/processes/process.go:3-6):

1. **Text Generation** (text2text.go) → Calls Ollama LLM API to generate story from prompt
2. **Speech Generation** (text2speech.go) → Invokes Python TTS script to convert text to audio
3. **Video Generation** (video.go) → Combines background video, audio, and captions

### Pipeline Execution
The `Generate` usecase (internal/usecases/generate.go:11-43) chains these processes together:
- Creates temporary working directory
- Links processes: textProcess → speechProcess → videoProcess
- Each process updates the `GenerationContext` and passes it to the next handler
- Returns final video path on success

### Key Components

**GenerationContext** (internal/processes/process.go:8-14): Shared state object passed through the pipeline containing:
- TempDir: Temporary working directory
- Prompt: Original user input
- Text: LLM-generated story
- SpeechFile: Generated audio file path
- GeneratedVideo: Final video output path

**Video Generation Pipeline** (internal/processes/video.go):
1. Selects random background video from `resources/videos/` (video.go:76-95)
2. Overlays generated audio onto background video using FFmpeg (video.go:97-114)
3. Crops video to match audio duration (video.go:116-135)
4. Runs Python caption script to generate and attach subtitles (video.go:154-170)

**Python Scripts** (scripts/):
- `tts.py`: Uses Facebook MMS-TTS model for text-to-speech conversion
- `captions.py`: Uses Whisper for transcription and MoviePy for subtitle rendering (max 18 words per caption)

### Docker Services

**docker-compose.yaml** defines two services:
- `ollama`: LLM service on port 11434 with persistent volume for model storage
- `worker`: Go application on port 8080, mounts `./generated` for output files

### Environment Variables

Required environment variables (set in docker-compose.yaml or .env):
- `PORT`: Server port (default: 8080)
- `ENV`: Environment mode (local/production)
- `OLLAMA_URL`: Ollama API endpoint
- `OLLAMA_MODEL`: Model name (e.g., orca-mini)

### Customization Points

**Background Videos**: Place MP4 files in `resources/videos/`. Video is randomly selected on each generation (video.go:76-95).

**Ollama Model**: Change `OLLAMA_MODEL` environment variable. Browse models at https://ollama.com/library.

**Caption Settings**: Adjust `max_words_per_caption` in scripts/captions.py:54 to control subtitle length.
