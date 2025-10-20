FROM python:3.13-slim

# Piper server docs:
# https://github.com/OHF-Voice/piper1-gpl/blob/main/docs/API_HTTP.md
RUN python3 -m pip install piper-tts[http]

RUN mkdir /models
VOLUME /models
ENV PIPER_MODEL=sv_SE-nst-medium.onnx

EXPOSE 5000

CMD python3 -m piper.http_server --model /models/${PIPER_MODEL} --data-dir /models --port 5000
