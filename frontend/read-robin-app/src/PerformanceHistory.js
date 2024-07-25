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
                contentID: quizDoc.data().contentID,
                title:
                  quizDoc.data().title ||
                  quizDoc.data().url ||
                  quizDoc.data().audio_url ||
                  quizDoc.data().video_url ||
                  quizDoc.data().pdf_url,
                url: quizDoc.data().url,
                audio_url: quizDoc.data().audio_url,
                video_url: quizDoc.data().video_url,
                pdf_url: quizDoc.data().pdf_url,
                attempts,
              };
            })
          );
          setQuizzes(quizzesData);
          console.log("Fetched quizzes data:", quizzesData);
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
      (quiz.attempts || []).filter((attempt) => {
        const attemptDate = new Date(attempt.createdAt.seconds * 1000);
        return now - attemptDate <= timeFrameMs;
      })
    );

    console.log("Filtered attempts for stats:", filteredAttempts);

    const totalQuizzes = filteredAttempts.length;
    const totalQuestions = filteredAttempts.reduce(
      (sum, attempt) =>
        sum + (attempt.responses ? attempt.responses.length : 0),
      0
    );
    const averageScore = filteredAttempts.length
      ? (
          filteredAttempts.reduce((sum, attempt) => sum + attempt.score, 0) /
          filteredAttempts.length
        ).toFixed(2)
      : 0;

    return { totalQuizzes, totalQuestions, averageScore };
  };

  const stats = calculateStats(timeFrame);
  console.log("Calculated stats:", stats);

  const getXLabels = (timeFrame) => {
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
      (quiz.attempts || []).filter((attempt) => {
        const attemptDate = new Date(attempt.createdAt.seconds * 1000);
        return now - attemptDate <= timeFrameMs;
      })
    );

    const xLabels = filteredAttempts.map((attempt) => {
      const attemptDate = new Date(attempt.createdAt.seconds * 1000);
      if (timeFrame === "24h") {
        return attemptDate.getHours() + ":00";
      } else if (timeFrame === "7d") {
        return attemptDate.toLocaleDateString("en-US", { weekday: "short" });
      } else {
        return attemptDate.toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        });
      }
    });

    console.log("X Labels:", xLabels);
    return xLabels;
  };

  const groupAttemptsByContentType = (attempts) => {
    const grouped = {
      URL: [],
      PDF: [],
      Audio: [],
      Video: [],
    };

    attempts.forEach((attempt) => {
      if (attempt.url) {
        grouped.URL.push(attempt);
      } else if (attempt.pdf_url) {
        grouped.PDF.push(attempt);
      } else if (attempt.audio_url) {
        grouped.Audio.push(attempt);
      } else if (attempt.video_url) {
        grouped.Video.push(attempt);
      }
    });

    console.log("Grouped attempts by content type:", grouped);
    return grouped;
  };

  const getScoresByContentType = (timeFrame) => {
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
      (quiz.attempts || []).filter((attempt) => {
        const attemptDate = new Date(attempt.createdAt.seconds * 1000);
        return now - attemptDate <= timeFrameMs;
      })
    );

    console.log("Filtered attempts for scores:", filteredAttempts);

    const groupedAttempts = groupAttemptsByContentType(filteredAttempts);

    return {
      URL: groupedAttempts.URL.map((attempt) => ({
        x: new Date(attempt.createdAt.seconds * 1000),
        y: attempt.score,
      })),
      PDF: groupedAttempts.PDF.map((attempt) => ({
        x: new Date(attempt.createdAt.seconds * 1000),
        y: attempt.score,
      })),
      Audio: groupedAttempts.Audio.map((attempt) => ({
        x: new Date(attempt.createdAt.seconds * 1000),
        y: attempt.score,
      })),
      Video: groupedAttempts.Video.map((attempt) => ({
        x: new Date(attempt.createdAt.seconds * 1000),
        y: attempt.score,
      })),
    };
  };

  const scoresByContentType = getScoresByContentType(timeFrame);
  console.log("Scores by content type:", scoresByContentType);

  const scoresOverTime = {
    datasets: [
      {
        label: "URL",
        data: scoresByContentType.URL,
        fill: false,
        backgroundColor: "rgba(75,192,192,1)",
        borderColor: "rgba(75,192,192,1)",
        borderWidth: 1,
        parsing: {
          xAxisKey: "x",
          yAxisKey: "y",
        },
      },
      {
        label: "PDF",
        data: scoresByContentType.PDF,
        fill: false,
        backgroundColor: "rgba(54,162,235,1)",
        borderColor: "rgba(54,162,235,1)",
        borderWidth: 1,
        parsing: {
          xAxisKey: "x",
          yAxisKey: "y",
        },
      },
      {
        label: "Audio",
        data: scoresByContentType.Audio,
        fill: false,
        backgroundColor: "rgba(255,206,86,1)",
        borderColor: "rgba(255,206,86,1)",
        borderWidth: 1,
        parsing: {
          xAxisKey: "x",
          yAxisKey: "y",
        },
      },
      {
        label: "Video",
        data: scoresByContentType.Video,
        fill: false,
        backgroundColor: "rgba(153,102,255,1)",
        borderColor: "rgba(153,102,255,1)",
        borderWidth: 1,
        parsing: {
          xAxisKey: "x",
          yAxisKey: "y",
        },
      },
    ],
  };

  console.log("Scores over time:", scoresOverTime);

  const contentTypes = {
    labels: ["URL", "PDF", "Audio", "Video"],
    datasets: [
      {
        label: "Content Types",
        data: quizzes.reduce(
          (acc, quiz) => {
            if (quiz.url) acc[0]++;
            else if (quiz.pdf_url) acc[1]++;
            else if (quiz.audio_url) acc[2]++;
            else if (quiz.video_url) acc[3]++;
            return acc;
          },
          [0, 0, 0, 0]
        ),
        backgroundColor: [
          "rgba(75,192,192,0.2)",
          "rgba(54,162,235,0.2)",
          "rgba(255,206,86,0.2)",
          "rgba(153,102,255,0.2)",
        ],
        borderColor: [
          "rgba(75,192,192,1)",
          "rgba(54,162,235,1)",
          "rgba(255,206,86,1)",
          "rgba(153,102,255,1)",
        ],
        borderWidth: 1,
      },
    ],
  };

  console.log("Content types:", contentTypes);

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
            options={{
              plugins: { legend: { display: true, position: "bottom" } },
              scales: {
                x: { title: { display: true, text: "Time" } },
                y: { title: { display: true, text: "Score" } },
              },
            }}
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
