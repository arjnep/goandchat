      const chooseContainer = document.getElementById("choose-container");
      const joinContainer = document.getElementById("join-container");
      const createContainer = document.getElementById("create-container");
      const messageContainer = document.getElementById("message-container");
      const card = document.getElementById("card");

      let chatArea = document.getElementById("chat-area");
      let socket;
      let username = "";
      let roomCode = "";

      initialTransition();

      function initializeWebSocket() {
        socket = new WebSocket(
          "ws://" + window.location.host + "/ws?code=" + roomCode
        );
        socket.onopen = function () {
          console.log("WebSocket connection established.");
        };
        socket.onmessage = function (event) {
          console.log("On message");
          let messageData = JSON.parse(event.data);
          console.log(messageData);
          displayMessage(messageData);
        };
        socket.onclose = function (event) {
          console.log("WebSocket connection closed:", event);
          if (!event.wasClean) {
            alert(
              "Connection closed unexpectedly. Please try reloading the page."
            );
            location.reload();
          } else {
            toaster("success", "Disconnected !");
            goback();
            chatArea.innerHTML = "";
          }
        };
        socket.onerror = function (event) {
          console.error("WebSocket error observed:", event);
        };
      }

      function joinChatroom() {
        const usernameInput = document.getElementById("username-input");
        const roomCodeInput = document.getElementById("room-code-input");
        username = usernameInput.value.trim();
        roomCode = roomCodeInput.value.trim();
        if (username === "" || roomCode === "") {
          toaster("error", "! Username and Room Code cannot be empty.");
          return;
        }

        fetch(`/ws?code=${roomCode}`)
          .then((response) => {
            if (response.status === 404) {
              toaster("error", "! Chatroom not found");
              return;
            }
            initializeWebSocket();
            document.getElementById("choose-container").classList.add("hidden");
            document.getElementById("join-container").classList.add("hidden");
            document.getElementById("create-container").classList.add("hidden");
            document
              .getElementById("message-container")
              .classList.remove("hidden");
            toaster("success", "Connected !");
            usernameInput.value = "";
            roomCodeInput.value = "";
          })
          .catch(() => {
            toaster("error", "Failed to connect to chatroom.");
          });
      }

      function createChatroom() {
        const usernameInput = createContainer.querySelector("#username-input");
        username = usernameInput.value.trim();
        if (username === "") {
          alert("Username cannot be empty.");
          return;
        }
        fetch("/create_chatroom", { method: "POST" })
          .then((response) => response.json())
          .then((data) => {
            roomCode = data.code;
            initializeWebSocket();
            document.getElementById("choose-container").classList.add("hidden");
            document.getElementById("join-container").classList.add("hidden");
            document.getElementById("create-container").classList.add("hidden");
            document
              .getElementById("message-container")
              .classList.remove("hidden");
            document.getElementById(
              "chatroom-code"
            ).textContent = `${roomCode}`;
            toaster("success", "Chatroom created!");
            usernameInput.value = "";
          });
      }

      function sendMessage() {
        let chatInput = document.getElementById("message-input");
        let message = chatInput.value;
        if (message === "") {
          return;
        }
        let messageData = {
          sender: username,
          content: message,
          time: Date.now(),
        };
        socket.send(JSON.stringify(messageData));
        chatInput.value = "";
      }

      function disconnect() {
        if (socket) {
          socket.close();
        }
      }

      function displayMessage(messageData) {
        if (
          messageData.sender === "system" &&
          messageData.content.startsWith("User count:")
        ) {
          let userCount = messageData.content.split(":")[1].trim();
          document.getElementById(
            "chatroom-users"
          ).textContent = `${userCount}`;
        } else if (
          messageData.sender === "system" &&
          messageData.content.startsWith("Room code:")
        ) {
          document.getElementById(
            "chatroom-code"
          ).textContent = `${messageData.content.split(":")[1].trim()}`;
        } else {
          let chatElement = document.createElement("div");
          chatElement.classList.add("chat");
          if (messageData.sender === username) {
            chatElement.classList.add("chat-go");
            chatElement.innerHTML = `<div class="chat-header"><span id="user">${messageData.sender}</span></div>`;
            chatElement.innerHTML += `<div class="chat-bubble">${messageData.content}</div>`;
          } else {
            chatElement.classList.add("chat-come");
            chatElement.innerHTML = `<div class="chat-header"><span id="user">${messageData.sender}</span></div>`;
            chatElement.innerHTML += `<div class="chat-bubble">${messageData.content}</div>`;
          }
          chatArea.appendChild(chatElement);
          chatArea.scrollTop = chatArea.scrollHeight;
        }
      }

      function goback() {
        setTimeout(function () {
          toggleTransition();
          chooseContainer.classList.remove("hidden");
          messageContainer.classList.add("hidden");
          createContainer.classList.add("hidden");
          joinContainer.classList.add("hidden");
        }, 1000);
        toggleTransition();
      }

      function createChatroomContainer() {
        setTimeout(function () {
          toggleTransition();
          createContainer.classList.remove("hidden");
          chooseContainer.classList.add("hidden");
          joinContainer.classList.add("hidden");
        }, 1000);
        toggleTransition();
      }

      function joinChatroomContainer() {
        setTimeout(function () {
          toggleTransition();
          joinContainer.classList.remove("hidden");
          chooseContainer.classList.add("hidden");
          createContainer.classList.add("hidden");
        }, 1000);
        toggleTransition();
      }

      function toaster(type, message) {
        const toasterDiv = document.getElementById("toaster");
        const messageSpan = toasterDiv.querySelector(".message");

        toasterDiv.classList.remove("success", "error");

        switch (type) {
          case "success":
            toasterDiv.classList.add("success");
            break;
          case "error":
            toasterDiv.classList.add("error");
            break;
          default:
            toasterDiv.classList.add("success");
        }

        messageSpan.textContent = message;

        toasterDiv.style.top = "4%";
        setTimeout(function () {
          toasterDiv.style.top = "-100px";
        }, 3000);
      }

      function initialTransition() {
        setTimeout(function () {
          card.classList.add("active");
        }, 100);
      }

      function toggleTransition() {
        card.classList.toggle("active");
      }

      window.onbeforeunload = function () {
        if (socket) {
          socket.close();
        }
      };