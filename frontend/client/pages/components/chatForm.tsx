import React, { useState, useEffect } from "react";

interface ws {
  connection: string;
  nameServer: string;
}

const ChatForm: React.FC<ws> = ({ connection, nameServer }) => {
  const [message, setMessage] = useState("");
  const websocket = new WebSocket(`ws://localhost:8888/ws/${connection}`);
  const [messageList, setMessageList] = useState<string[]>([]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (websocket.readyState === WebSocket.OPEN) {
      websocket.send(message);
      setMessage("");
    } else {
      console.error("WebSocket connection is not open.");
    }
  };

  useEffect(() => {
    // Add a message event listener to the WebSocket
    websocket.addEventListener("message", (event) => {
      const receivedMessage = event.data;
      // Update the messageList state with the new message
      setMessageList((prevMessageList) => [...prevMessageList, receivedMessage]);
    });

    // Clean up the event listener when the component unmounts
    return () => {
      websocket.removeEventListener("message", (event) => {
        // Clean up any resources if needed
      });
    };
  }, []); // Add websocket as a dependency to ensure correct cleanup

  return (
    <div className="flex flex-col items-center justify-center h-screen text-black">
      <h2 className="text-white text-2xl text-center">{nameServer}</h2>
      <div className="w-96 h-96 bg-gray-200 rounded-lg p-4">
        <div className="h-64 overflow-y-scroll mb-4">
          {messageList.map((message, index) => (
            <h2 key={index} className="text-black text-left gap-5">
              {message}
            </h2>
          ))}
        </div>
        <form onSubmit={handleSubmit} className="mb-4">
          <input
            type="text"
            placeholder="Type your message..."
            className="w-full p-2 border rounded-md"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
          <button
            type="submit"
            className="mt-2 p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring focus:ring-blue-300"
          >
            Send
          </button>
        </form>
      </div>
    </div>
  );
};

export default ChatForm;
