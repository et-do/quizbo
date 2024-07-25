import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import { Bar, Line } from "react-chartjs-2";
import { Chart as ChartJS, registerables } from "chart.js";
import "./PerformanceHistory.css";

ChartJS.register(...registerables);

function PerformanceHistory({
  user,
  activePersona,
  setPage,
  setAttemptID,
  setContentID,
}) {
  const [quizzes, setQuizzes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [timeFrame, setTimeFrame] = useState("24h");

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

  const calculateStats = (timeFrame) => {
    const now = new Date();
    let timeFrameMs;

    switch (timeFrame) {
      case "24h":
        timeFrameMs = 24 * 60 * 60 * 1000;
        break;
      case "7d":
        timeFrameMs = 7 * 24 * 60 * 60 * 1000;
        break;
      case "1m":
        timeFrameMs = 30 * 24 * 60 * 60 * 1000;
        break;
      default:
        timeFrameMs = 24 * 60 * 60 * 1000;
    }

    const filteredAttempts = quizzes.flatMap((quiz) =>
      quiz.attempts.filter((attempt) => {
        const attemptDate = new Date(attempt.createdAt.seconds * 1000);
        return now - attemptDate <= timeFrameMs;
      })
    );

    const totalQuizzes = new Set(
      filteredAttempts.map((attempt) => attempt.contentID)
    ).size;
    const totalQuestions = filteredAttempts.length;
    const averageScore = filteredAttempts.length
      ? (
          filteredAttempts.reduce((sum, attempt) => sum + attempt.score, 0) /
          filteredAttempts.length
        ).toFixed(2)
      : 0;

    return { totalQuizzes, totalQuestions, averageScore };
  };

  const stats = calculateStats(timeFrame);

  const scoresOverTime = {
    labels: quizzes.flatMap((quiz) =>
      quiz.attempts
        .filter((attempt) => {
          const attemptDate = new Date(attempt.createdAt.seconds * 1000);
          return (
            new Date() - attemptDate <=
            (timeFrame === "24h"
              ? 24 * 60 * 60 * 1000
              : timeFrame === "7d"
              ? 7 * 24 * 60 * 60 * 1000
              : 30 * 24 * 60 * 60 * 1000)
          );
        })
        .map((attempt) =>
          new Date(attempt.createdAt.seconds * 1000).toLocaleDateString()
        )
    ),
    datasets: [
      {
        label: "Scores",
        data: quizzes.flatMap((quiz) =>
          quiz.attempts
            .filter((attempt) => {
              const attemptDate = new Date(attempt.createdAt.seconds * 1000);
              return (
                new Date() - attemptDate <=
                (timeFrame === "24h"
                  ? 24 * 60 * 60 * 1000
                  : timeFrame === "7d"
                  ? 7 * 24 * 60 * 60 * 1000
                  : 30 * 24 * 60 * 60 * 1000)
              );
            })
            .map((attempt) => attempt.score)
        ),
        fill: false,
        backgroundColor: "rgba(75,192,192,1)",
        borderColor: "rgba(75,192,192,1)",
        borderWidth: 1,
      },
    ],
  };

  const contentTypes = {
    labels: ["URL", "PDF", "Audio", "Video"],
    datasets: [
      {
        label: "Content Types",
        data: quizzes.reduce(
          (acc, quiz) => {
            if (quiz.title.includes("http")) acc[0]++;
            else if (quiz.title.endsWith(".pdf")) acc[1]++;
            else if (quiz.title.endsWith(".mp3")) acc[2]++;
            else acc[3]++;
            return acc;
          },
          [0, 0, 0, 0]
        ),
        backgroundColor: [
          "rgba(255,99,132,0.2)",
          "rgba(54,162,235,0.2)",
          "rgba(255,206,86,0.2)",
          "rgba(75,192,192,0.2)",
        ],
        borderColor: [
          "rgba(255,99,132,1)",
          "rgba(54,162,235,1)",
          "rgba(255,206,86,1)",
          "rgba(75,192,192,1)",
        ],
        borderWidth: 1,
      },
    ],
  };

  return (
    <div className="performance-history">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Performance History</h2>
      <div className="timeframe-select">
        <label>
          <input
            type="radio"
            name="timeframe"
            value="24h"
            checked={timeFrame === "24h"}
            onChange={(e) => setTimeFrame(e.target.value)}
          />
          <span>Last 24 Hours</span>
        </label>
        <label>
          <input
            type="radio"
            name="timeframe"
            value="7d"
            checked={timeFrame === "7d"}
            onChange={(e) => setTimeFrame(e.target.value)}
          />
          <span>Last 7 Days</span>
        </label>
        <label>
          <input
            type="radio"
            name="timeframe"
            value="1m"
            checked={timeFrame === "1m"}
            onChange={(e) => setTimeFrame(e.target.value)}
          />
          <span>Last 1 Month</span>
        </label>
        <div className="slider"></div>
      </div>
      <div className="stats-container">
        <div className="stats-card">
          <div className="stat">
            <h3>Quizzes Taken</h3>
            <p>{stats.totalQuizzes}</p>
          </div>
          <div className="stat">
            <h3>Questions Answered</h3>
            <p>{stats.totalQuestions}</p>
          </div>
          <div className="stat">
            <h3>Average Score</h3>
            <p>{stats.averageScore}%</p>
          </div>
        </div>
      </div>
      <div className="charts-container">
        <div className="chart-card">
          <h3>Scores Over Time</h3>
          <Line
            data={scoresOverTime}
            options={{ plugins: { legend: { display: false } } }}
          />
        </div>
        <div className="chart-card">
          <h3>Content Types</h3>
          <Bar
            data={contentTypes}
            options={{ plugins: { legend: { display: false } } }}
          />
        </div>
      </div>
      {quizzes.length > 0 && (
        <div className="quizzes-list">
          <h3>Quizzes and Attempts</h3>
          {quizzes.map((quiz) => (
            <div key={quiz.id} className="quiz-card">
              <h4>{quiz.title}</h4>
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
