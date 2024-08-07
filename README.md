# Quizbo
![Quizbo Logo](frontend/read-robin-app/src/logo.png)

**Turn Anything into a Quiz**: Your Personalized Knowledge Platform for Articles, Research Papers, Podcasts, Videos, and more!

## Overview
Quizbo is a versatile application designed to transform any media into interactive quizzes. 

Whether it's text from URLs, video, audio, or images, Quizbo generates quizzes to help users engage with content and enhance their learning.

Users can store their quiz history, compare their performance with peers, and continuously improve their knowledge.

## Features
- **Multi-Media Quiz Generation**: Generate quizzes from text, video, audio, and images.
- **User Authentication**: Secure login with Google authentication.
- **Quiz History**: Track and review past quizzes.
- **Peer Comparison**: Compare quiz results with peers.
- **Interactive UI**: Engaging and user-friendly interface.

## Repository Layout
```bash
.
├── backend/
│   ├── .devcontainer/     # Development container configuration
│   ├── handlers/          # HTTP handler functions
│   ├── models/            # Common custom types
│   ├── middleware/        # Logging functions
│   ├── services/          # External Service interaction (Gemini, Firestore)
│   ├── utils/             # Utility functions
│   ├── main.go            # Entry point of the backend application
│   └── go.mod             # Go module file
└── frontend/
    ├── .devcontainer/     # Development container configuration
    ├── src/               # React source files
    ├── public/            # Public assets
    ├── package.json       # Node.js dependencies and scripts
    └── firebase.json      # Firebase configuration
```

#### Attributions

Logo: <a href="https://www.flaticon.com/free-icons/robin" title="robin icons">Robin icons created by Freepik - Flaticon</a>