import React, { useEffect, useState } from "react";
import "./IntroScreen.css";
import logo from "./logo.png"; // Import the logo

const slides = [
  {
    title: "Welcome to ReadRobin!",
    text: "Your AI Companion for Smarter Comprehension",
  },
  {
    title: "Generate Quizzes",
    text: "Turn your content into quizzes for improved learning",
  },
  {
    title: "Track Your Progress",
    text: "See your quiz history and improve over time",
  },
];

const IntroScreen = ({ onFinish }) => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const [fadeOut, setFadeOut] = useState(false);

  useEffect(() => {
    if (currentSlide < slides.length) {
      const timer = setTimeout(() => {
        setCurrentSlide((prev) => prev + 1);
      }, 3000); // Adjust the timeout as needed
      return () => clearTimeout(timer);
    } else {
      const fadeTimer = setTimeout(() => {
        setFadeOut(true);
        const finishTimer = setTimeout(onFinish, 1000); // Adjust as needed for fade-out duration
        return () => clearTimeout(finishTimer);
      }, 1000);
      return () => clearTimeout(fadeTimer);
    }
  }, [currentSlide, onFinish]);

  return (
    <div className={`intro-screen ${fadeOut ? "fade-out" : ""}`}>
      <img src={logo} alt="Logo" className="intro-logo" />
      <div className="slides-container">
        {slides.map((slide, index) => (
          <div
            key={index}
            className={`slide ${index === currentSlide ? "active" : ""} ${
              index < currentSlide ? "left" : ""
            }`}
          >
            <h2>{slide.title}</h2>
            <p>{slide.text}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default IntroScreen;
