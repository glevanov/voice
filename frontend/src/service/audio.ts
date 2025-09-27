let currentAudio: string | null = null;

export const play = async (
  audioElement: HTMLAudioElement | null,
  audioPath: string,
) => {
  try {
    const response = await fetch(`http://localhost:3002/${audioPath}`);
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
