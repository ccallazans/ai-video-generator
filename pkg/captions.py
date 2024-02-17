import logging
import sys
import moviepy.editor as mp
from moviepy.video.tools.subtitles import SubtitlesClip
import os
import speech_recognition as sr
import textwrap

def parse_video(video_path, name, temp_folder):
    logging.basicConfig(level = logging.INFO)
    logging.info('Starting parse_video')
    output_path = os.path.join(temp_folder, name)
    generate_subtitles(video_path, output_path, temp_folder, name)
    return output_path

def generate_subtitles(video_path, output_path, temp_folder, name):
    logging.info('Starting generate_subtitles')
    try:
        audio_path = f"{temp_folder}/temp_audio_{name}.wav"
        extract_audio(video_path, audio_path)
        subtitle_text = transcribe_audio(audio_path)
        generate_srt(subtitle_text, output_path)

        video = mp.VideoFileClip(video_path)
        text_generator = _create_text_generator(video)

        logging.info('Starting SubtitlesClip')
        subtitles = SubtitlesClip(output_path, text_generator)

        if subtitles.duration > 0:  # Check for non-empty subtitles
            video = mp.CompositeVideoClip([video, subtitles.set_position('center')])

        logging.info('Starting write_videofile')
        video.write_videofile(output_path, codec='libx264', audio_codec='aac')
    finally:
        os.remove(audio_path)  # Clean up audio file

def extract_audio(video_path, audio_path):
    logging.info('Starting extract_audio')
    mp.VideoFileClip(video_path).audio.write_audiofile(audio_path)

def transcribe_audio(audio_path):
    logging.info('Starting transcribe_audio')
    recognizer = sr.Recognizer()
    with sr.AudioFile(audio_path) as source:
        try:
            audio_data = recognizer.record(source)
            return recognizer.recognize_google(audio_data)
        except sr.UnknownValueError:
            return ""

def generate_srt(subtitle_text: str, output_srt_path: str) -> None:
    logging.info('Starting generate_srt')
    chunks = chunk_text(subtitle_text, 40)
    time_interval = 2.48 # Higher is slower
    srt_content = ''
    for i, chunk in enumerate(chunks):
        start_time = i * time_interval
        end_time = start_time + time_interval
        start_time_formatted = "{:02.0f}:{:02.0f}:{:02.0f},000".format(start_time // 3600, (start_time % 3600) // 60, start_time % 60)
        end_time_formatted = "{:02.0f}:{:02.0f}:{:02.0f},000".format(end_time // 3600, (end_time % 3600) // 60, end_time % 60)
        srt_content += f"{i+1}\n{start_time_formatted} --> {end_time_formatted}\n{chunk}\n\n"
    with open(output_srt_path, 'w') as f:
        f.write(srt_content)

def _create_text_generator(video):
    logging.info('Starting _create_text_generator')
    max_chars_per_line = 18
    font_path = "./pkg/config/bold_font.ttf"
    fontsize = int(video.size[1] * 0.05)
    color = "#FFFF00"
    stroke_color = "black"
    stroke_width = 5
    return lambda txt: mp.TextClip(
        '\n'.join(textwrap.fill(line, width=max_chars_per_line) for line in txt.split('\n')),
        font=font_path, fontsize=fontsize, color=color, stroke_color=stroke_color, stroke_width=stroke_width
    )

def chunk_text(text, max_length):
    logging.info('Starting chunk_text')
    chunks = []
    current_chunk = ""
    for word in text.split():
        if len(current_chunk) + len(word) + 1 > max_length:
            chunks.append(current_chunk.rstrip())
            current_chunk = ""
        current_chunk += word + " "
    return chunks + [current_chunk.rstrip()]
    

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python main.py <parameter>")
        sys.exit(1)

    video_to_process_path = temp_folder = sys.argv[1]
    final_video_path = sys.argv[2]
    temp_folder = sys.argv[3]
    parse_video(video_to_process_path, final_video_path, temp_folder)