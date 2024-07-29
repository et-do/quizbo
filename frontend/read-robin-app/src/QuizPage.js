import React, { useState, useEffect } from "react";
import "./QuizPage.css";
import { db } from "./firebase";
import { doc, setDoc, getDoc, Timestamp } from "firebase/firestore";

function QuizPage({ user, activePersona, setPage, contentID, quizID }) {
  const [questions, setQuestions] = useState([]);
  const [responses, setResponses] = useState({});
  const [status, setStatus] = useState({});
  const [explanations, setExplanations] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [attemptID, setAttemptID] = useState(() => {
    const timestamp = new Date().toISOString();
    return `${quizID}@${timestamp}`;
  });
  const [quizTitle, setQuizTitle] = useState("");
  const [showPopup, setShowPopup] = useState(false);
  const [popupContent, setPopupContent] = useState("");
  const [submitting, setSubmitting] = useState({});

  useEffect(() => {
    const fetchQuizData = async () => {
      try {
        if (!user || !activePersona || !activePersona.id) {
          throw new Error("User or active persona is not defined");
        }

        const quizRef = doc(
          db,
          "users",
          user.uid,
          "personas",
          activePersona.id,
          "quizzes",
          contentID
        );
        const quizDoc = await getDoc(quizRef);
        if (quizDoc.exists()) {
          setQuizTitle(quizDoc.data().title || ""); // Fetch title
        }

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
    fetchQuizData();
  }, [contentID, quizID, user, activePersona]);

  const handleResponseChange = (e, index) => {
    const newResponses = { ...responses, [index]: e.target.value };
    setResponses(newResponses);
  };

  const calculateScore = (responses) => {
    const totalQuestions = responses.length;
    const correctAnswers = responses.filter(
      (response) => response.status === "Correct"
    ).length;
    return (correctAnswers / totalQuestions) * 100;
  };

  const handleSubmitResponse = async (index, questionID) => {
    const userResponse = responses[index];
    const questionData = questions[index];

    if (status[index] !== undefined) {
      alert("You have already submitted a response for this question.");
      return;
    }

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
      persona: {
        id: activePersona.id,
        name: activePersona.name,
        role: activePersona.role,
        language: activePersona.language,
        difficulty: activePersona.difficulty,
      },
    };

    try {
      setSubmitting({ ...submitting, [index]: true });
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
        [index]: data.status.trim() === "PASS" ? "Correct" : "Incorrect",
      };
      const newExplanations = {
        ...explanations,
        [index]: data.explanation,
      };
      setStatus(newStatus);
      setExplanations(newExplanations);

      const attemptRef = doc(
        db,
        "users",
        user.uid,
        "personas",
        activePersona.id,
        "quizzes",
        contentID,
        "attempts",
        attemptID
      );
      const attemptDoc = await getDoc(attemptRef);
      let existingResponses = [];

      if (attemptDoc.exists()) {
        existingResponses = attemptDoc.data().responses || [];
      }

      const updatedResponses = [
        ...existingResponses,
        {
          questionID: questionID,
          question: questionData.question,
          answer: questionData.answer,
          reference: questionData.reference,
          userResponse: userResponse,
          status: data.status.trim() === "PASS" ? "Correct" : "Incorrect",
        },
      ];

      const score = calculateScore(updatedResponses);

      await setDoc(
        attemptRef,
        {
          attemptID: attemptID,
          createdAt: Timestamp.now(),
          responses: updatedResponses,
          score: score,
          title: quizTitle,
        },
        { merge: true }
      );
    } catch (error) {
      console.error("Error submitting response: ", error);
      if (error.code === "permission-denied") {
        console.error("Permission denied! Check your Firestore rules.");
      }
    } finally {
      setSubmitting({ ...submitting, [index]: false });
    }
  };

  const handleRetakeQuiz = () => {
    const timestamp = new Date().toISOString();
    setAttemptID(`${quizID}@${timestamp}`);
    setResponses({});
    setStatus({});
    setExplanations({});
    setSubmitting({});
  };

  const handleShowExplanation = (index) => {
    setPopupContent(explanations[index]);
    setShowPopup(true);
  };

  const handleClosePopup = () => {
    setShowPopup(false);
    setPopupContent("");
  };

  return (
    <div className="quiz-page">
      <button className="back-button" onClick={() => setPage("selection")}>
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
                  disabled={submitting[index]}
                >
                  Submit
                </button>
              </div>
              {submitting[index] && <div className="loading-spinner"></div>}
              {status[index] && (
                <div
                  className={
                    status[index] === "Correct" ? "correct" : "incorrect"
                  }
                >
                  {status[index]}
                  <button
                    className="explanation-button"
                    onClick={() => handleShowExplanation(index)}
                  >
                    ℹ️
                  </button>
                </div>
              )}
            </div>
          ))}
          <button className="retake-button" onClick={handleRetakeQuiz}>
            Retake Quiz
          </button>
        </div>
      )}
      {showPopup && (
        <div className="popup">
          <div className="popup-content">
            <span className="close-button" onClick={handleClosePopup}>
              &times;
            </span>
            <p>{popupContent}</p>
          </div>
        </div>
      )}
    </div>
  );
}

export default QuizPage;
