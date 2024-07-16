import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import "./Sidebar.css";

function Sidebar({ user, setContentID, setAttemptID, setPage }) {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);
  const [expandedUrl, setExpandedUrl] = useState(null);

  useEffect(() => {
    const fetchQuizzes = async () => {
      if (user) {
        try {
          const quizzesRef = collection(db, "users", user.uid, "quizzes");
          const querySnapshot = await getDocs(quizzesRef);
          const quizzesData = await Promise.all(
            querySnapshot.docs.map(async (quizDoc) => {
              const attemptsRef = collection(
                db,
                "users",
                user.uid,
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
                url: quizDoc.data().url,
                attempts,
              };
            })
          );
          console.log("Fetched quizzes with attempts:", quizzesData);
          setQuizzes(quizzesData);
          setLoading(false);
        } catch (error) {
          console.error("Error fetching quizzes:", error);
          setLoading(false);
        }
      }
    };
    fetchQuizzes();
  }, [user]);

  const handleAttemptClick = (contentID, attemptID) => {
    setContentID(contentID);
    setAttemptID(attemptID);
    setPage("attemptPage");
  };

  const toggleSidebar = () => {
    setIsOpen(!isOpen);
  };

  const toggleUrl = (url) => {
    setExpandedUrl(expandedUrl === url ? null : url);
  };

  if (loading) {
    return <div className="sidebar">Loading...</div>;
  }

  return (
    <>
      <button
        className={`sidebar-toggle ${isOpen ? "open" : ""}`}
        onClick={toggleSidebar}
      >
        {isOpen ? ">" : "<"}
      </button>
      <div className={`sidebar ${isOpen ? "open" : ""}`}>
        <h2>Quiz History</h2>
        <ul>
          {quizzes.map((quiz) => (
            <li key={quiz.id}>
              <h3 onClick={() => toggleUrl(quiz.url)}>{quiz.url}</h3>
              {expandedUrl === quiz.url && (
                <ul>
                  {quiz.attempts.length > 0 ? (
                    quiz.attempts.map((attempt) => (
                      <li
                        key={attempt.attemptID}
                        onClick={() =>
                          handleAttemptClick(quiz.id, attempt.attemptID)
                        }
                      >
                        {attempt.createdAt
                          ? new Date(
                              attempt.createdAt.seconds * 1000
                            ).toLocaleString()
                          : "N/A"}
                        <br />
                        {attempt.score ? `${attempt.score}%` : "N/A"}
                      </li>
                    ))
                  ) : (
                    <li>No attempts found</li>
                  )}
                </ul>
              )}
            </li>
          ))}
        </ul>
      </div>
      <div
        className={`overlay ${isOpen ? "open" : ""}`}
        onClick={toggleSidebar}
      ></div>
    </>
  );
}

export default Sidebar;
