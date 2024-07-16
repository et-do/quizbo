import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { doc, getDoc } from "firebase/firestore";
import "./AttemptPage.css";

function AttemptPage({ user, contentID, attemptID, setPage }) {
  const [attempt, setAttempt] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quizUrl, setQuizUrl] = useState("");

  useEffect(() => {
    const fetchAttempt = async () => {
      if (user && contentID && attemptID) {
        const quizRef = doc(db, "users", user.uid, "quizzes", contentID);
        const quizDoc = await getDoc(quizRef);
        if (quizDoc.exists()) {
          setQuizUrl(quizDoc.data().url);
        }

        const attemptRef = doc(
          db,
          "users",
          user.uid,
          "quizzes",
          contentID,
          "attempts",
          attemptID
        );
        const attemptDoc = await getDoc(attemptRef);
        if (attemptDoc.exists()) {
          setAttempt(attemptDoc.data());
        } else {
          console.error("No such document!");
        }
        setLoading(false);
      }
    };
    fetchAttempt();
  }, [user, contentID, attemptID]);

  const getScoreClass = (score) => {
    if (score <= 50) return "red";
    if (score >= 80) return "green";
    return "";
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!attempt) {
    return <div>No data found</div>;
  }

  return (
    <div className="attempt-page">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>{quizUrl}</h2>
      <p className="score-text">
        Score:{" "}
        <span className={`score-value ${getScoreClass(attempt.score)}`}>
          {attempt.score}%
        </span>
      </p>
      <ul>
        {attempt.responses.map((response, index) => (
          <li key={index}>
            <div className="response-item">
              <span className="response-title">Question:</span>
              <span>{response.question}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Correct Answer:</span>
              <span>{response.answer}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Your Response:</span>
              <span>{response.userResponse}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Status:</span>
              <span>{response.status}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Reference:</span>
              <span>{response.reference}</span>
            </div>
            {index < attempt.responses.length - 1 && <hr />}
          </li>
        ))}
      </ul>
      <p className="attempt-id">Attempt ID: {attempt.attemptID}</p>
    </div>
  );
}

export default AttemptPage;
