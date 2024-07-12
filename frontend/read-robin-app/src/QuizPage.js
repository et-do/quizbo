import React, { useState, useEffect } from "react";
import "./App.css";
import { db } from "./firebase";
import {
  doc,
  setDoc,
  updateDoc,
  arrayUnion,
  Timestamp,
} from "firebase/firestore";

function QuizPage({ user, setPage, contentID, quizID }) {
  const [questions, setQuestions] = useState([]);
  const [responses, setResponses] = useState({});
  const [status, setStatus] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [attemptID, setAttemptID] = useState(() => {
    const timestamp = new Date().toISOString();
    return `${quizID}_${timestamp}`;
  });

  useEffect(() => {
    const fetchQuestions = async () => {
      try {
        const res = await fetch(
          `https://read-robin-dev-6yudia4zva-nn.a.run.app/quiz/${contentID}/${quizID}`
        );
        const data = await res.json();
        if (data.questions) {
          setQuestions(data.questions);
        } else {
          setError("Error fetching questions");
        }
        setLoading(false);
      } catch (error) {
        setError("Error fetching questions");
        setLoading(false);
      }
    };
    fetchQuestions();
  }, [contentID, quizID]);

  const handleResponseChange = (e, index) => {
    const newResponses = { ...responses, [index]: e.target.value };
    setResponses(newResponses);
  };

  const handleSubmitResponse = async (index, questionID) => {
    const userResponse = responses[index];
    const questionData = questions[index];

    if (status[index] !== undefined) {
      alert("You have already submitted a response for this question.");
      return;
    }

    // Validate that all necessary fields are defined
    if (
      userResponse === undefined ||
      questionData.question === undefined ||
      questionData.answer === undefined ||
      questionData.reference === undefined ||
      questionID === undefined
    ) {
      console.error("One or more fields are undefined", {
        userResponse,
        question: questionData.question,
        answer: questionData.answer,
        reference: questionData.reference,
        questionID,
      });
      return;
    }

    const payload = {
      content_id: contentID,
      quiz_id: quizID,
      question_id: questionID,
      user_response: userResponse,
    };

    try {
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/submit-response`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        }
      );

      if (!res.ok) {
        throw new Error(`Error submitting response: ${res.statusText}`);
      }

      const data = await res.json();
      const newStatus = {
        ...status,
        [index]: data.status === "PASS" ? "Correct" : "Incorrect",
      };
      setStatus(newStatus);

      // Save the user's response to Firestore along with the correct answer, question, and reference
      const attemptRef = doc(
        db,
        "users",
        user.uid,
        "quizzes",
        contentID,
        "attempts",
        attemptID
      );
      await setDoc(
        attemptRef,
        {
          attemptID: attemptID,
          createdAt: Timestamp.now(),
          responses: arrayUnion({
            questionID: questionID,
            question: questionData.question,
            answer: questionData.answer,
            reference: questionData.reference,
            userResponse: userResponse,
            status: data.status,
          }),
        },
        { merge: true }
      );
    } catch (error) {
      console.error("Error submitting response: ", error);
      if (error.code === "permission-denied") {
        console.error("Permission denied! Check your Firestore rules.");
      }
    }
  };

  const handleRetakeQuiz = () => {
    const timestamp = new Date().toISOString();
    setAttemptID(`${quizID}_${timestamp}`); // Generate a new attempt ID
    setResponses({});
    setStatus({});
  };

  return (
    <div className="quiz-page">
      <button className="back-button" onClick={() => setPage("quizForm")}>
        Back
      </button>
      <h2>Quiz</h2>
      {error && !loading && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="loading-spinner"></div>}
      {!loading && questions.length > 0 && (
        <div>
          {questions.map((item, index) => (
            <div key={index} className="quiz-item">
              <h3>Question {index + 1}</h3>
              <p>{item.question}</p>
              <div className="response-container">
                <input
                  type="text"
                  value={responses[index] || ""}
                  onChange={(e) => handleResponseChange(e, index)}
                />
                <button
                  onClick={() => handleSubmitResponse(index, item.question_id)}
                >
                  Submit
                </button>
              </div>
              {status[index] && (
                <div
                  className={
                    status[index] === "Correct" ? "correct" : "incorrect"
                  }
                >
                  {status[index]}
                </div>
              )}
            </div>
          ))}
          <button className="retake-button" onClick={handleRetakeQuiz}>
            Retake Quiz
          </button>
        </div>
      )}
    </div>
  );
}

export default QuizPage;
