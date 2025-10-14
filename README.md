# Voice Assistant
A real-time voice assistant with speech synthesis using Piper TTS and LLM integration.

## Binaries setup
You need to have Piper and Whisper installed:
- Piper: https://github.com/OHF-Voice/piper1-gpl
- Whisper: https://github.com/ggml-org/whisper.cpp

You'll also need models for both:
- Piper: https://huggingface.co/rhasspy/piper-voices/tree/main
- Whisper: https://huggingface.co/ggerganov/whisper.cpp/tree/main

Place models in the `models` directory.

You'll also need `ffmpeg` for audio conversion.

To set up Piper and download swedish models run `sh bin/setup_piper.sh`.
The script assumes you are on Linux amd64.

Whisper you'll have to compile yourself and add it to `PATH`,
as there is no universal binary.

My models folder looks like this:
```shell
# $ tree models
models
├── en_US-amy-medium.onnx
├── en_US-amy-medium.onnx.json
├── ggml-base.bin
├── ggml-large-v3.bin
├── ggml-medium.bin
├── sv_SE-nst-medium.onnx
└── sv_SE-nst-medium.onnx.json
```

## License
MIT