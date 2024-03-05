const chatbox = document.getElementById("chatbox");
const messageInput = document.getElementById("messageInput");
const sendButton = document.getElementById("sendButton");

const socket = new WebSocket("ws://localhost:8080/ws"); // Change the URL to your server's WebSocket URL

socket.onopen = function () {
    console.log("Connected to the server.");
};

socket.onmessage = function(event) {
    console.log("Data received: " + event.data);

    // Распарсить JSON данные от сервера
    var message = JSON.parse(event.data);
    var author = data.author;
    var message = data.body;

    // Вызов функции displayMessage с полученными данными
    displayMessage(author, message);
};

socket.onerror = function(error) {
    console.log("Error " + error.message);
};

socket.onclose = function(event) {
    if (event.wasClean) {
        console.log('Connection close successfull');
    } else {
        console.log('Connection lost'); // например, "убит" процесс сервера
    }
    console.log('Code: ' + event.code + ' reason: ' + event.reason + ' error: ' + event.error + ' event: ' + event);
};

sendButton.addEventListener("click", function () {
    const message = messageInput.value;
    if (message.trim() !== "") {
        sendMessage(message);
        messageInput.value = "";
    }
});

// Обработчик клика на кнопку Send
document.getElementById("sendButton").addEventListener("click", function(event) {
    event.preventDefault();
    sendMessage();
});

function sendMessage(message) {
    const messageObj = {
        "author": "Admin",
        "body": messageInput.value
    };
    socket.send(JSON.stringify(messageObj));
}

function displayMessage(message) {
    const messageElement = document.createElement("div");
    messageElement.classList.add("message");
    messageElement.innerHTML = `<span class="user">${message.author}:</span> ${message.body}`;
    chatbox.appendChild(messageElement);
    chatbox.scrollTop = chatbox.scrollHeight;
}