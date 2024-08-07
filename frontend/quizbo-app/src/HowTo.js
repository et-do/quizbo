import React from "react";
import "./App.css";

function HowTo({ setPage }) {
  return (
    <div className="howto-page">
      <h1>How to Use Quizbo</h1>
      <p>
        Quizbo is your AI companion for smarter comprehension. Follow the steps
        below to get the most out of Quizbo!
      </p>
      <h2>1. Manage Personas</h2>
      <p>
        Go to{" "}
        <strong
          onClick={() => setPage("personas")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Personas
        </strong>{" "}
        to create a persona. Creating a persona will allow you to take tests
        however you'd like! Define who you are, your difficulty choice, language
        choice, and give that persona a name. Now all your quizzes under that
        active persona will be tailored for you! To change or add personas, go
        to the{" "}
        <strong
          onClick={() => setPage("personas")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Personas
        </strong>{" "}
        tab.
      </p>
      <h2>2. Generate a Quiz</h2>
      <p>
        Once you have a persona, go and{" "}
        <strong
          onClick={() => setPage("selection")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          generate a quiz
        </strong>
        ! You can generate quizzes from website HTML through URLs, PDFs, audio,
        or video, as well as plain text. Once you've generated a quiz, you can
        answer as few questions as you'd like, and get instant feedback from
        Quizbo on how you were right or wrong, as well as learning tips!
      </p>
      <h2>3. Performance History</h2>
      <p>
        Once you've completed an attempt, your performance history will be saved
        on the sidebar and in the{" "}
        <strong
          onClick={() => setPage("performanceHistory")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Performance History
        </strong>{" "}
        tab. Clicking on any of the attempts in the sidebar will bring you to
        the attempts page where you can see your answers, the correct answers,
        and references from the text itself. Check out the{" "}
        <strong
          onClick={() => setPage("performanceHistory")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Performance History
        </strong>{" "}
        tab to get historic trends on your learning performance.
      </p>
      <h2>4. Manage Content</h2>
      <p>
        Want to quiz yourself again on the same content? Navigate to the{" "}
        <strong
          onClick={() => setPage("contentManagement")}
          style={{ cursor: "pointer", color: "#4169e1" }}
        >
          Manage Content
        </strong>{" "}
        tab to regenerate quizzes based on content you've already submitted. You
        can also see your latest scores to ensure you are up to date with all
        your important content.
      </p>
    </div>
  );
}

export default HowTo;
