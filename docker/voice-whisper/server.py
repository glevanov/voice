#!/usr/bin/env python3
"""
Simple HTTP server wrapper for whisper-cli.
Accepts requests to transcribe audio files and returns the text.
"""

import os
import subprocess
import json
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

WHISPER_MODEL = os.environ.get('WHISPER_MODEL', 'ggml-large-v3.bin')
LANGUAGE = os.environ.get('LANGUAGE', 'Swedish')
MODEL_PATH = f'/models/{WHISPER_MODEL}'
PORT = 3000


class WhisperHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/transcribe':
            self._transcribe_audio()
        else:
            self.send_error(404, 'Not Found')

    def _transcribe_audio(self):
        audio_file = '/audio/question.wav'

        if not os.path.exists(audio_file):
            self.send_error(404, f'Audio file not found: {audio_file}')
            return

        try:
            cmd = [
                'whisper-cli',
                '--no-prints',
                '--no-timestamps',
                '--language', LANGUAGE,
                '--model', MODEL_PATH,
                '--file', audio_file
            ]
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=300
            )

            if result.returncode != 0:
                self.send_error(500, f'Whisper error: {result.stderr}')
                return

            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()

            self.wfile.write(result.stdout.strip().encode())

        except subprocess.TimeoutExpired:
            self.send_error(504, 'Transcription timeout')
        except Exception as e:
            self.send_error(500, f'Internal server error: {str(e)}')

def main():
    server_address = ('', PORT)
    httpd = HTTPServer(server_address, WhisperHandler)
    print(f'Starting Whisper HTTP server on port {PORT}')
    httpd.serve_forever()


if __name__ == '__main__':
    main()

