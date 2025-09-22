<script>
    import "./app.css";
    import { get } from "svelte/store";
    import {
        messages,
        addUserMessage,
        addAssistantMessage,
        clearMessages,
    } from "./store/messages.js";

    let text = "";
    let response = "";
    let audio = null;
    let audioElement = null;
    let connectionStatus = "Connecting...";

    const ws = new WebSocket("ws://localhost:3002/ws");

    const handleSend = () => {
        if (ws.readyState === WebSocket.OPEN) {
            const userMessage = text;
            addUserMessage(userMessage);

            // Get current messages (which now includes the new user message)
            const currentMessages = get(messages);

            // Send complete chat history
            const payload = {
                messages: currentMessages,
            };
            ws.send(JSON.stringify(payload));
            text = "";
        } else {
            alert("WebSocket connection is not open");
        }
    };

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            response = data.text;
            addAssistantMessage(data.text);

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
            addAssistantMessage(event.data);
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

    <div class="chat-history">
        <div class="chat-header">
            <h2>Chat History</h2>
            <button on:click={clearMessages} class="clear-btn"
                >Clear History</button
            >
        </div>
        <div class="messages">
            {#each $messages as message}
                <div class="message {message.role}">
                    <div class="role">
                        {message.role === "user" ? "You" : "Assistant"}:
                    </div>
                    <div class="content">{message.content}</div>
                </div>
            {/each}
            {#if $messages.length === 0}
                <div class="no-messages">No messages yet...</div>
            {/if}
        </div>
    </div>

    <div class="input-section">
        <textarea bind:value={text} placeholder="Type your message here..."
        ></textarea>
        <button
            on:click={handleSend}
            disabled={connectionStatus !== "Connected"}>Send</button
        >
    </div>

    {#if response}
        <h2>Latest Response:</h2>
        <div class="response">{response}</div>
    {/if}

    {#if audio}
        <audio bind:this={audioElement} controls src={audio} />
    {/if}
</div>

<style>
    .chat-history {
        margin: 20px 0;
        border: 1px solid #ccc;
        border-radius: 8px;
        max-height: 400px;
        overflow-y: auto;
    }

    .chat-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 10px 15px;
        border-bottom: 1px solid #eee;
        background-color: #f8f9fa;
    }

    .chat-header h2 {
        margin: 0;
        font-size: 1.2em;
    }

    .clear-btn {
        padding: 5px 10px;
        background-color: #dc3545;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.9em;
    }

    .clear-btn:hover {
        background-color: #c82333;
    }

    .messages {
        padding: 15px;
    }

    .message {
        margin-bottom: 15px;
        padding: 10px;
        border-radius: 8px;
    }

    .message.user {
        background-color: #e3f2fd;
        margin-left: 20px;
    }

    .message.assistant {
        background-color: #f1f8e9;
        margin-right: 20px;
    }

    .role {
        font-weight: bold;
        margin-bottom: 5px;
        color: #555;
    }

    .content {
        white-space: pre-wrap;
    }

    .no-messages {
        text-align: center;
        color: #666;
        font-style: italic;
        padding: 20px;
    }

    .input-section {
        margin: 20px 0;
    }

    .input-section textarea {
        width: 100%;
        margin-bottom: 10px;
    }
</style>
