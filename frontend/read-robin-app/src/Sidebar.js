import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import "./Sidebar.css";

function Sidebar({ user, setContentID, setAttemptID, setPage }) {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    const fetchQuizzes = async () => {
      if (user) {
        try {
          const quizzesRef = collection(db, "users", user.uid, "quizzes");
          const querySnapshot = await getDocs(quizzesRef);
          const quizzesData = querySnapshot.docs.map((doc) => ({
            id: doc.id,
            ...doc.data(),
          }));
          console.log("Fetched quizzes:", quizzesData); // Debugging statement
          setQuizzes(quizzesData);
          setLoading(false);
        } catch (error) {
          console.error("Error fetching quizzes:", error);
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
              <h3>{quiz.title}</h3>
              {quiz.attempts ? (
                quiz.attempts.map((attempt) => (
                  <p
                    key={attempt.attemptID}
                    onClick={() =>
                      handleAttemptClick(quiz.id, attempt.attemptID)
                    }
                  >
                    Attempt on:{" "}
                    {attempt.createdAt
                      ? new Date(
                          attempt.createdAt.seconds * 1000
                        ).toLocaleDateString()
                      : "N/A"}
                    <br />
                    Score: {attempt.score ? `${attempt.score}%` : "N/A"}
                  </p>
                ))
              ) : (
                <p>No attempts found</p>
              )}
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}

export default Sidebar;
