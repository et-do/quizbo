import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import "./PerformanceHistory.css";

function PerformanceHistory({
  user,
  activePersona,
  setPage,
  setAttemptID,
  setContentID,
}) {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchQuizzes = async () => {
      if (user && activePersona) {
        try {
          const quizzesRef = collection(
            db,
            "users",
            user.uid,
            "personas",
            activePersona.id,
            "quizzes"
          );
          const querySnapshot = await getDocs(quizzesRef);
          const quizzesData = await Promise.all(
            querySnapshot.docs.map(async (quizDoc) => {
              const attemptsRef = collection(
                db,
                "users",
                user.uid,
                "personas",
                activePersona.id,
                "quizzes",
                quizDoc.id,
                "attempts"
              );
              const attemptsSnapshot = await getDocs(attemptsRef);
              const attempts = attemptsSnapshot.docs
                .map((attemptDoc) => ({
                  attemptID: attemptDoc.id,
                  ...attemptDoc.data(),
                }))
                .sort((a, b) => b.createdAt.seconds - a.createdAt.seconds); // Sort by recency
              return {
                id: quizDoc.id,
                title: quizDoc.data().title || quizDoc.data().url,
                attempts,
              };
            })
          );
          setQuizzes(quizzesData);
          setLoading(false);
        } catch (error) {
          console.error("Error fetching quizzes:", error);
          setLoading(false);
        }
      } else {
        setLoading(false);
      }
    };
    fetchQuizzes();
  }, [user, activePersona]);

  const handleAttemptClick = (contentID, attemptID) => {
    setContentID(contentID);
    setAttemptID(attemptID);
    setPage("attemptPage");
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="performance-history">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Performance History</h2>
      {quizzes.length === 0 ? (
        <p>No quizzes found.</p>
      ) : (
        <div className="dashboard">
          {quizzes.map((quiz) => (
            <div key={quiz.id} className="quiz-card">
              <h3>{quiz.title}</h3>
              {quiz.attempts.length > 0 ? (
                <ul>
                  {quiz.attempts.map((attempt) => (
                    <li
                      key={attempt.attemptID}
                      onClick={() =>
                        handleAttemptClick(quiz.id, attempt.attemptID)
                      }
                      className="attempt-item"
                    >
                      <div className="attempt-info">
                        <p>
                          <strong>Date:</strong>{" "}
                          {attempt.createdAt
                            ? new Date(
                                attempt.createdAt.seconds * 1000
                              ).toLocaleString()
                            : "N/A"}
                        </p>
                        <p>
                          <strong>Score:</strong>{" "}
                          {attempt.score ? `${attempt.score}%` : "N/A"}
                        </p>
                      </div>
                      <button className="view-details-button">
                        View Details
                      </button>
                    </li>
                  ))}
                </ul>
              ) : (
                <p>No attempts found</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default PerformanceHistory;
