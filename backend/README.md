# Quizbo Backend

## Overview

This project is the backend part of the Quizbo application. It is a Go application running on Cloud Run integrated with Google Cloud Platform (GCP) services, including Firestore, Vertex AI (Gemini model), and Firebase for authentication. The backend handles content extraction, quiz generation, and storing/retrieving quiz data in Firestore.

## Local Development

To develop locally within the devcontainer:

1. **Open the project in a devcontainer**.
2. **Add account credentials** to the `secrets/` directory.

## Deploying Changes

To deploy changes to Cloud Run, follow these guidelines:

- **Develop Branch**: Merging your changes to the `develop` branch will automatically trigger a deployment to the Cloud Run development environment.
- **Main Branch**: Merging your changes to the `main` branch will automatically trigger a deployment to the production Cloud Run environment.

## Project Structure
```
.
├── .devcontainer/
├── handlers/ # Contains HTTP handler functions
│ ├── submit.go
│ └── quiz.go
├── models/ # Contains common custom types
│ └── firebase_collection_schemas.go
├── middleware/ # Contains logging functions
│ └── logging.go
├── services/ # Contains service files for interacting with external APIs and Firestore
│ ├── firestore.go
│ └── gemini.go
├── utils/ # Utility functions (e.g., fetching HTML content)
│ └── html_fetcher.go
├── secrets/ # Credential Keys
├── main.go # Entry point of the application, sets up routes and middleware
└── go.mod # Go module file
```

## API Endpoints

### 1. Submit URL

- **Endpoint**: `/submit`
- **Method**: POST
- **Description**: Submits a URL to extract content, generate a quiz, and save it to Firestore.
- **Request Body**:
    ```json
    {
        "url": "http://example.com"
    }
    ```
- **Response**:
    ```json
    {
        "status": "success",
        "url": "http://example.com",
        "content_id": "abcd1234",
        "quiz_id": "0001"
    }
    ```

### 2. Get Quiz by ContentID and QuizID
- **Endpoint**: `/quiz/{contentID}/{quizID}`
- **Method**: GET
- **Description**: Retrieves quiz questions from Firestore by ContentID and QuizID.
- **Response**:
    ```json
    {
        "questions": [
            {
                "question_id": "0001",
                "question": "What is the purpose of the example domain?",
                "answer": "The 'Example Domain' is for use in illustrative examples in documents.",
                "reference": "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission."
            },
            {
                "question_id": "0002",
                "question": "Where can you find more information about the example domain?",
                "answer": "You can find more information on the IANA website.",
                "reference": "More information can be found on the IANA website."
            }
        ]
    }
    ```

### 3. Submit Quiz Response

- **Endpoint**: `/submit-response`
- **Method**: POST
- **Description**: Submits a user's response to a quiz question for review.
- **Request Body**:
    ```json
    {
        "content_id": "abcd1234",
        "quiz_id": "0001",
        "question_id": "0001",
        "user_response": "It is used for examples in documents."
    }
    ```
- **Response**:
    ```json
    {
        "status": "PASS"
    }
    ```


## Testing
Test files are written alongside the files they are testing (I.e. "services/firestore.go", "services/firestore_test.go")
# Unit Tests
Run unit tests using the following command:
```sh
go test ./...
```

## Common Issues
Ensure you have set up your GCP credentials and project ID correctly.
If you encounter issues with Firestore, ensure the Firestore Emulator is running or you have proper access to the Firestore database.

## Emergency Procedures
In case of an emergency where you need to pause the backend service:

Stop the Google Cloud Run service:

```sh
gcloud run services update read-robin --no-traffic
```
Redeploy the service to restore traffic.

## Contributing
To contribute to this project:

- Fork the repository.
- Create a feature branch (git checkout -b feature-branch).
- Commit your changes (git commit -am 'Add new feature').
- Push to the branch (git push origin feature-branch).
- Create a new Pull Request.