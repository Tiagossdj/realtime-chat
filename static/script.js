const ws = new WebSocket('ws://localhost:8080/ws');

const chat = document.getElementById('chat')
const messageInput = document.getElementById('message')
const sendButton = document.getElementById('send')

// Recebe mensagens do servidor e exibe no chat
ws.onmessage = function (event){
  const chat = document.getElementById('chat');
  const message = document.createElement('div');
  message.textContent = event.data;
  chat.appendChild(message);
  chat.scrollTop = chat.scrollHeight;
};

//Envia mensagens para o servidor 
sendButton.addEventListener('click', () => {
  const message = messageInput.value;
  if (message) {
    ws.send(message);
    messageInput.value = '';
  }
});