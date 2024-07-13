import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { doc, getDoc } from "firebase/firestore";
import "./AttemptPage.css";

function AttemptPage({ user, contentID, attemptID, setPage }) {
  const [attempt, setAttempt] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAttempt = async () => {
      if (user && contentID && attemptID) {
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
      <h2>Attempt Details</h2>
      <p>Attempt ID: {attempt.attemptID}</p>
      <p>Score: {attempt.score}%</p>
      <ul>
        {attempt.responses.map((response, index) => (
          <li key={index}>
            <p>Question: {response.question}</p>
            <p>Correct Answer: {response.answer}</p>
            <p>Your Response: {response.userResponse}</p>
            <p>Status: {response.status}</p>
            <p>Reference: {response.reference}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default AttemptPage;
