let currentAudio: string | null = null;

const ANSWER_URL = "http://localhost:4003/answer.wav";

export const play = async (audioElement: HTMLAudioElement | null) => {
  try {
    const response = await fetch(ANSWER_URL);
    if (!response.ok) {
      throw new Error("Failed to fetch audio");
    }

    const blob = await response.blob();
    currentAudio = URL.createObjectURL(blob);

    if (audioElement) {
      audioElement.src = currentAudio;
      setTimeout(() => {
        audioElement.play().catch((error) => {
          console.error("Autoplay failed:", error);
        });
      }, 100);
    }

    return currentAudio;
  } catch (error) {
    console.error("Error fetching audio:", error);
    throw error;
  }
};

export const cleanup = () => {
  if (currentAudio) {
    URL.revokeObjectURL(currentAudio);
    currentAudio = null;
  }
};
