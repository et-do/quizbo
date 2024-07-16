// IntroScreen.js
import React, { useEffect } from "react";
import "./IntroScreen.css";

const IntroScreen = ({ onFinish }) => {
  useEffect(() => {
    const timer = setTimeout(onFinish, 5000); // Adjust the timeout as needed
    return () => clearTimeout(timer);
  }, [onFinish]);

  return (
    <div className="intro-screen">
      <div className="slide">
        <h2>Welcome to ReadRobin!</h2>
        <p>Your AI Companion for Smarter Comprehension</p>
      </div>
      <div className="slide">
        <h2>Generate Quizzes</h2>
        <p>Turn articles into quizzes for better learning</p>
      </div>
      <div className="slide">
        <h2>Track Your Progress</h2>
        <p>See your quiz history and improve over time</p>
      </div>
    </div>
  );
};

export default IntroScreen;
