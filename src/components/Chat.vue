<template>
  <div class="chat-container">
    <h2>Chat</h2>
    <div class="messages">
      <div
        v-for="(message, index) in messages"
        :key="index"
        :class="{'my-message': message.sender === currentUser, 'other-message': message.sender !== currentUser}"
      >
        <strong>{{ message.sender }}:</strong> {{ message.text }}
      </div>
    </div>
    <form @submit.prevent="sendMessage" class="chat-form">
      <input
        v-model="newMessage"
        type="text"
        placeholder="Type your message here..."
        required
      />
      <button type="submit">Send</button>
    </form>
  </div>
</template>

<script>
export default {
  name: 'Chat',
  data() {
    return {
      // In a real application, currentUser would come from your Vuex store.
      currentUser: 'Me',
      newMessage: '',
      messages: [
        { sender: 'Alice', text: 'Hello!' },
        { sender: 'Me', text: 'Hi, how are you?' },
      ],
    };
  },
  methods: {
    sendMessage() {
      if (this.newMessage.trim() !== '') {
        // Append the new message to the messages array.
        this.messages.push({
          sender: this.currentUser,
          text: this.newMessage.trim(),
        });
        // Clear the input field.
        this.newMessage = '';
        // In a real app, you would also send this message to your server here.
      }
    },
  },
};
</script>

<style scoped>
.chat-container {
  padding: 2rem;
}
.messages {
  border: 1px solid #ccc;
  padding: 1rem;
  height: 300px;
  overflow-y: auto;
  margin-bottom: 1rem;
}
.my-message {
  text-align: right;
  margin: 0.5rem 0;
  background-color: #e1f5fe;
  padding: 0.5rem;
  border-radius: 4px;
}
.other-message {
  text-align: left;
  margin: 0.5rem 0;
  background-color: #f1f1f1;
  padding: 0.5rem;
  border-radius: 4px;
}
.chat-form {
  display: flex;
}
.chat-form input {
  flex: 1;
  padding: 0.5rem;
  font-size: 1rem;
}
.chat-form button {
  padding: 0.5rem 1rem;
  margin-left: 0.5rem;
  font-size: 1rem;
}
</style>

