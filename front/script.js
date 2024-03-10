const chatbox = document.getElementById("chatbox");
const messageInput = document.getElementById("messageInput");
const sendButton = document.getElementById("sendButton");

const socket = new WebSocket("ws://localhost:8080/chat/ws");

socket.onopen = function () {
    console.log("Connected to the server.");
};

socket.onmessage = function(event) {
    console.log("Data received: " + event.data);

    // Распарсить JSON данные от сервера
    var message = JSON.parse(event.data);
    var author = message.author;
    var body = message.body;

    // Вызов функции displayMessage с полученными данными
    displayMessage(author, body);
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

messageInput.addEventListener("keypress", function(event){
    var key = event.key 
    if (key === "Enter"){
        sendButton.click()
    }
})

function sendMessage(message) {
    const messageObj = {
        "author": "Admin",
        "body": messageInput.value
    };
    socket.send(JSON.stringify(messageObj));
}

function displayMessage(author, body) {
    const messageElement = document.createElement("div");
    messageElement.classList.add("message");
    messageElement.innerHTML = `<span class="user">${author}:</span> ${body}`;

    console.log("displayed: " + author + " " + body)

    chatbox.appendChild(messageElement);
    chatbox.scrollTop = chatbox.scrollHeight;
}