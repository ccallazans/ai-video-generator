#!/usr/bin/env python3
# Credit goes to
# https://github.com/pacifio/autocap
# github.com/pacifio

import os
import argparse
import subprocess
from datetime import timedelta
import whisper
from moviepy.editor import VideoFileClip, CompositeVideoClip, TextClip
from moviepy.video.tools.subtitles import SubtitlesClip
import string
import random

VALID_MODES = ("attach", "generate")
TEMP_FILE = "temp.mp3"
OUTPUT_SRT = "output.srt"
OUTPUT_VID = "./generated"

class VideoManager:
    def __init__(self, path: str) -> None:
        self.path = path
        self.video = VideoFileClip(path)
        self.audio_path = os.path.join(os.path.dirname(path), TEMP_FILE)
        self.extract_audio()

    def extract_audio(self) -> None:
        if self.video.audio is not None:
            self.video.audio.write_audiofile(self.audio_path, codec="mp3")
        else:
            print("video has no audio, quitting")
            exit()

class Utility:
    def __init__(self, path: str) -> None:
        self.path = path

    def file_exists(self) -> bool:
        return os.path.exists(self.path)

class SubtitleGenerator:
    def __init__(self, videomanager: VideoManager, random_word: str) -> None:
        self.videomanager = videomanager
        self.directory = os.path.dirname(self.videomanager.path)
        self.srt_path = os.path.join(self.directory, OUTPUT_SRT)
        self.output_vid_path = os.path.join(OUTPUT_VID, random_word)

    def generate(self) -> None:
        model = whisper.load_model("base")
        transcribe = model.transcribe(audio=self.videomanager.audio_path, fp16=False)
        segments = transcribe["segments"]

        max_words_per_caption = 18  # Adjust this value as needed

        with open(self.srt_path, "w", encoding="utf-8") as f:
            segment_id = 1
            for seg in segments:
                start = str(0) + str(timedelta(seconds=int(seg["start"]))) + ",000"
                end = str(0) + str(timedelta(seconds=int(seg["end"]))) + ",000"
                text = seg["text"].strip()
                
                words = text.split()
                for i in range(0, len(words), max_words_per_caption):
                    text_chunk = " ".join(words[i:i + max_words_per_caption])
                    segment = f"{segment_id}\n{start} --> {end}\n{text_chunk}\n\n"
                    f.write(segment)
                    segment_id += 1
        print("subtitles generated")

    def attach(self) -> None:
        self.generate()
        if os.path.exists(self.srt_path):
            subtitles = SubtitlesClip(
                self.srt_path,
                lambda txt: TextClip(
                    txt,
                    font="FreeSans-Bold",
                    fontsize=60,
                    color="white",
                    method='caption',
                    size=(self.videomanager.video.w - 100, None),
                    align="center"
                )
            )

            video_with_subtitles = CompositeVideoClip([
                self.videomanager.video,
                subtitles.set_position(('center', 0.4), relative=True)
            ])

            video_with_subtitles.write_videofile(self.output_vid_path, codec="libx264")
            print(f"saved to {self.output_vid_path}")


def check_ffmpeg() -> bool:
    try:
        result = subprocess.run(['ffmpeg', '-version'], capture_output=True, text=True)
        return result.returncode == 0 and 'ffmpeg' in result.stdout
    except FileNotFoundError:
        return False

def generate_random_word(length=7):
    letters = string.ascii_letters
    return ''.join(random.choice(letters) for _ in range(length))

def main() -> None:
    parser = argparse.ArgumentParser(description="auto caption generator v1.0")
    parser.add_argument("mode", metavar="mode", type=str, help="operation mode (attach|generate)")
    parser.add_argument("path", metavar="path", type=str, help="filepath of the video")
    parser.add_argument("random_word", metavar="random_word", type=str, help="random word for the output filename")
    args = parser.parse_args()

    if not check_ffmpeg():
        print("ffmpeg must be installed to run this script, quitting")
        exit()

    mode = args.mode
    path = args.path
    random_word = args.random_word

    if mode in VALID_MODES and os.path.exists(path):
        videomanager = VideoManager(path)
        subtitle_generator = SubtitleGenerator(videomanager, random_word)

        if mode == "attach":
            subtitle_generator.attach()
        elif mode == "generate":
            subtitle_generator.generate()
    else:
        print("invalid mode or file path, quitting")

if __name__ == "__main__":
    main()
