# Quizbo
<a href="https://quizbo.app">
  <img src="frontend/quizbo-app/src/logo.png" alt="Quizbo Logo" width="300" height="auto">
</a>

**Turn Anything into a Quiz**: Your Personalized Knowledge Platform for Articles, Research Papers, Podcasts, Videos, and more!


## Overview
Quizbo is a versatile application designed to transform any media into interactive quizzes. 

Whether it's text from URLs, video, audio, or images, Quizbo generates quizzes to help users engage with content and enhance their learning.

Users can store their quiz history, compare their performance with peers, and continuously improve their knowledge.

## Features
- **Multi-Media Quiz Generation**: Generate quizzes from text, video, audio, and images.
- **Quiz History**: Track and review past quizzes. See scores over time and what type of content you quiz yourself on.
- **Tailored Peronas**: Create a user persona to tailor the quizzes to your specific needs. Control the audience, difficulty, and language of your quizzes
- **Question Insights**: Get helpful insights on each and every question, right or wrong.

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


## Fair Use & Data Privacy
- **Fair Use**: Any content uploaded by the user is presumed to be owned by the user and/or the user has the explicit permission to generate quizzes with it.
- **Data Privacy**: None of the uploaded data from users is accessed by others and can be permanently removed from the databases by each user at any point.
