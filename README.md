# AI Video Story Generator

This project utilizes Ollama LLM to create videos text ideas based on user prompts. It takes a prompt from the user and generates a story, which is then converted into audio using AI. Subsequently, captions are generated from the audio, and finally, these captions are merged into a video.

## How It Works

1. **Prompt Input**: Users provide a prompt to the system.
2. **Story Generation**: AI generates a story based on the provided prompt.
3. **Audio Generation**: The story is converted into audio using AI text-to-speech technology.
4. **Caption Generation**: AI generates captions from the audio.
5. **Video Creation**: Captions are merged into a video.

## How to Run

### Prerequisites
- Golang 1.21.1
- Docker installed on your system

### Installation
1. Clone the repository:
```
git clone https://github.com/your-username/ai-video-story-generator.git
cd ai-video-story-generator
```
2. Build the project:
```
make build
```
3. Start the Docker containers:
```
sudo docker-compose up -d && sudo docker exec -it $(sudo docker ps -aqf "name=ollama") ollama pull orca-mini
```

Once the containers are up and running, you can access the AI video story generator through its provided interface or API endpoints.

### Usage


####  Send post request with the content you want to be generated
```
curl --location 'localhost:1323/api/v1/generate' \
--header 'Content-Type: application/json' \
--data '{
    "message": "Tell me a story about Bahia"
}'
```

## Additional Notes
- To add more background videos, put it inside resources/videos folder. The video is random selected, so change on code if want.
- Ensure that Docker is properly configured and running on your system.
- Adjust any necessary configurations in the Dockerfile or docker-compose.yml file according to your requirements.
- For production use, ensure proper security measures and scalability considerations are implemented.

Feel free to explore and modify the project according to your needs! If you encounter any issues or have suggestions for improvement, please don't hesitate to open an issue or submit a pull request.
