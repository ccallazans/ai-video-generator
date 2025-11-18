import asyncio
import edge_tts
import sys

# Available voices:
# English (US): en-US-AriaNeural (female), en-US-GuyNeural (male)
# Spanish (Spain): es-ES-ElviraNeural (female), es-ES-AlvaroNeural (male)
# Spanish (Mexico): es-MX-DaliaNeural (female), es-MX-JorgeNeural (male)
# Portuguese (Brazil): pt-BR-FranciscaNeural (female), pt-BR-AntonioNeural (male)
# More voices: https://speech.microsoft.com/portal/voicegallery

DEFAULT_VOICE = "es-ES-ElviraNeural"

async def text_to_speech(text, output_filename, voice=DEFAULT_VOICE):
    communicate = edge_tts.Communicate(text, voice)
    await communicate.save(output_filename)
    print(f"Audio file written successfully to {output_filename} using voice {voice}")

if __name__ == "__main__":
    if len(sys.argv) < 3 or len(sys.argv) > 4:
        print("Usage: tts.py <text> <output_filename> [voice]")
        sys.exit(1)

    text = sys.argv[1]
    output_filename = sys.argv[2]
    voice = sys.argv[3] if len(sys.argv) == 4 else DEFAULT_VOICE

    asyncio.run(text_to_speech(text, output_filename, voice))