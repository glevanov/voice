<script>
    import "./app.css";

    let text = "";
    let response = "";
    let audio = null;
    let audioElement = null;
    let connectionStatus = "Connecting...";

    const ws = new WebSocket("ws://localhost:3002/ws");

    const handleSend = () => {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(text);
            text = ""; // Clear input after sending
        } else {
            alert("WebSocket connection is not open");
        }
    };

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            response = data.text;

            if (data.audio) {
                fetch(`http://localhost:3002/${data.audio}`)
                    .then((res) => {
                        if (res.ok) {
                            return res.blob();
                        }
                        throw new Error("Failed to fetch audio");
                    })
                    .then((blob) => {
                        audio = URL.createObjectURL(blob);
                        // Autoplay the audio after a short delay to ensure element is ready
                        setTimeout(() => {
                            if (audioElement) {
                                audioElement.play().catch((error) => {
                                    console.error("Autoplay failed:", error);
                                });
                            }
                        }, 100);
                    })
                    .catch((error) => {
                        console.error("Error fetching audio:", error);
                    });
            }
        } catch (error) {
            console.error("Error parsing response:", error);
            response = event.data; // Show raw response if JSON parsing fails
        }
    };

    ws.onopen = () => {
        connectionStatus = "Connected";
    };

    ws.onclose = () => {
        connectionStatus = "Disconnected";
    };

    ws.onerror = (error) => {
        connectionStatus = "Error";
        console.error("WebSocket error:", error);
    };
</script>

<div class="app">
    <h1>Voice Assistant</h1>
    <div class="status">Status: {connectionStatus}</div>
    <textarea bind:value={text} placeholder="Type your message here..."
    ></textarea>
    <button on:click={handleSend} disabled={connectionStatus !== "Connected"}
        >Send</button
    >
    <h2>Response:</h2>
    <div class="response">{response || "No response yet..."}</div>
    {#if audio}
        <audio bind:this={audioElement} controls src={audio} />
    {/if}
</div>
