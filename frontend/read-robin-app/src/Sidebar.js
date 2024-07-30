import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import "./Sidebar.css";

function Sidebar({ user, activePersona, setContentID, setAttemptID, setPage }) {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);
  const [expandedQuiz, setExpandedQuiz] = useState(null);

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
                score: parseFloat(attemptDoc.data().score), // Convert score to float
              }))
              .sort((a, b) => b.createdAt.seconds - a.createdAt.seconds); // Sort by recency
            return {
              id: quizDoc.id,
              title: quizDoc.data().title,
              content_type: quizDoc.data().content_type,
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

  useEffect(() => {
    fetchQuizzes();
  }, [user, activePersona]);

  const handleAttemptClick = (contentID, attemptID) => {
    setContentID(contentID);
    setAttemptID(attemptID);
    setPage("attemptPage");
  };

  const toggleSidebar = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      fetchQuizzes();
    }
  };

  const toggleQuiz = (quizID) => {
    setExpandedQuiz(expandedQuiz === quizID ? null : quizID);
  };

  const groupAttemptsByContentType = (quizzes) => {
    const grouped = {
      URL: [],
      PDF: [],
      Audio: [],
      Video: [],
    };

    quizzes.forEach((quiz) => {
      if (quiz.content_type === "URL") {
        grouped.URL.push(quiz);
      } else if (quiz.content_type === "PDF") {
        grouped.PDF.push(quiz);
      } else if (quiz.content_type === "Audio") {
        grouped.Audio.push(quiz);
      } else if (quiz.content_type === "Video") {
        grouped.Video.push(quiz);
      }
    });

    return grouped;
  };

  const groupedQuizzes = groupAttemptsByContentType(quizzes);

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
        {activePersona && <h3>Persona: {activePersona.name}</h3>}
        {quizzes.length === 0 ? (
          <p>No quizzes found.</p>
        ) : (
          <>
            {Object.entries(groupedQuizzes).map(([contentType, quizzes]) => (
              <div key={contentType}>
                <h3>{contentType}</h3>
                <ul>
                  {quizzes.map((quiz) => (
                    <li key={quiz.id} className="quiz-item">
                      <div
                        className="quiz-title"
                        onClick={() => toggleQuiz(quiz.id)}
                      >
                        {quiz.title}
                      </div>
                      {expandedQuiz === quiz.id && (
                        <ul className="attempts-list">
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
                                {attempt.score ? `${attempt.score}%` : "0.00%"}
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
            ))}
          </>
        )}
      </div>
      <div
        className={`overlay ${isOpen ? "open" : ""}`}
        onClick={toggleSidebar}
      ></div>
    </>
  );
}

export default Sidebar;
