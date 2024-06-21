# ChitChat

ChitChat is a simple, yet powerful chatroom application built using Go for the backend and basic HTML, CSS, and JavaScript for the frontend. Users can create chatrooms, share unique codes to invite others, and communicate in real-time. The final version is under development.

## Features

- **User Authentication:** Users can sign up and log in with their credentials stored in a JSON file.
- **Chatroom Creation:** Users can create new chatrooms that generate unique alphanumeric codes.
- **Join Chatrooms:** Users can join existing chatrooms using the unique code.
- **Real-Time Messaging:** Messages are sent and displayed in real-time with the latest message always visible without needing to scroll.
- **Responsive Design:** The application has a modern and responsive design, ensuring a great user experience on both desktop and mobile devices.

## Technologies Used

- **Go:** The backend server is built using Go, handling user authentication, chatroom management, and message handling.
- **HTML:** The structure of the web pages is built using HTML.
- **CSS:** Styling is done using CSS to create a modern and responsive design.
- **JavaScript:** Client-side scripting for handling form submissions and auto-scrolling in the chat window.
- **JSON:** User credentials are stored in a JSON file for simplicity.

## Usage

- **Sign Up:** Create a new account by providing a username and password.
- **Log In:** Log in with your credentials.
- **Create a Chatroom:** Click on the "Create Room" button to generate a new chatroom with a unique code.
- **Join a Chatroom:** Enter a valid chatroom code to join an existing chatroom.
- **Send Messages:** Type your message and click "Send" to communicate with others in the chatroom.

## Future updates

Some planned future updates are under development! These include -
- **Error Handling and security:** Implementing proper error handling for cases like duplicate chatrooms and improving application security.
- **User Profile Management:** Allowing users to create and update their profiles and change passwords.
- **Chatroom Management:** Implementing features like chatroom moderation and kicking users.
- **Robust Database:** Integrating MongoDB to store user credentials.
- **Redis for storing messages:** To improve performance and security, and bring it one-step closer to being a production-grade application.