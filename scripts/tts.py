from transformers import VitsModel, AutoTokenizer
import torch
import scipy.io.wavfile
import sys

def text_to_speech(text, output_filename):
    model = VitsModel.from_pretrained("facebook/mms-tts-eng")
    tokenizer = AutoTokenizer.from_pretrained("facebook/mms-tts-eng")

    inputs = tokenizer(text, return_tensors="pt")

    with torch.no_grad():
        output = model(**inputs).waveform

    waveform = output.float().numpy()
    sampling_rate = model.config.sampling_rate

    if waveform.ndim > 1:
        waveform = waveform.squeeze()

    if waveform.max() > 1.0 or waveform.min() < -1.0:
        waveform = waveform / max(abs(waveform.max()), abs(waveform.min()))

    if not (0 < sampling_rate <= 65535):
        raise ValueError(f"Invalid sampling rate: {sampling_rate}")

    scipy.io.wavfile.write(output_filename, rate=sampling_rate, data=waveform)
    print("WAV file written successfully to", output_filename)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit(1)
    
    text = sys.argv[1]
    output_filename = sys.argv[2]
    text_to_speech(text, output_filename)