#!/bin/bash

# Download and extract Piper binary
wget https://github.com/rhasspy/piper/releases/download/v1.2.0/piper_amd64.tar.gz
tar -xzf piper_amd64.tar.gz

# Download Swedish voice model (NST medium)
# More models available at:
# https://huggingface.co/rhasspy/piper-voices/tree/main
wget -P ../models https://huggingface.co/rhasspy/piper-voices/resolve/main/sv/sv_SE/nst/medium/sv_SE-nst-medium.onnx
wget -P ../models https://huggingface.co/rhasspy/piper-voices/resolve/main/sv/sv_SE/nst/medium/sv_SE-nst-medium.onnx.json
