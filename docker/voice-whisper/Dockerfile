# This is for buildind whisper.cpp with Vulkan support.
# The have ready to use CPU and CUDA docker images, but not for Vulkan.
# https://github.com/ggml-org/whisper.cpp?tab=readme-ov-file#docker
# Using Fedora here as Debian does not have glslc package out of the box
FROM fedora:42 AS builder

RUN dnf install -y \
    git \
    cmake \
    gcc \
    gcc-c++ \
    make \
    vulkan-headers \
    vulkan-loader-devel \
    vulkan-tools \
    glslang \
    glslc \
    mesa-vulkan-drivers \
    && dnf clean all

RUN git clone https://github.com/ggml-org/whisper.cpp.git
WORKDIR /whisper.cpp
RUN cmake -B build -DGGML_VULKAN=1
RUN cmake --build build -j --config Release

FROM python:3.13-slim

RUN apt update \
    && apt install -y \
    libvulkan1 \
    mesa-vulkan-drivers \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /whisper.cpp /whisper.cpp

ENV PATH="/whisper.cpp/build/bin:${PATH}"
ENV WHISPER_MODEL=ggml-large-v3.bin
ENV LANGUAGE=Swedish

RUN mkdir /models
VOLUME /models
RUN mkdir /audio
VOLUME /audio

RUN mkdir /server
COPY server.py /server/server.py

EXPOSE 3000

WORKDIR /server
ENTRYPOINT []

CMD ["python3", "server.py"]
