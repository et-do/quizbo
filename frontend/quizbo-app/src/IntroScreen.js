import React, { useState } from "react";
import "./IntroScreen.css";
import logo from "./logo.png";

const slides = [
  {
    title: "👋 Welcome to Quizbo",
    text: "Your AI Companion for Smarter, Personalized Comprehension. Powered by Google's Gemini.",
  },
  {
    title: "📝 Generate Quizzes from Any Content",
    text: "Turn any content—websites, PDFs, podcasts, and videos—into quizzes.",
  },
  {
    title: "💡 Tailored Learning Experience",
    text: "Whether you're a student, professional, or lifelong learner, generate quizzes that match your role, preferred language, and difficulty level",
  },
  {
    title: "🚀 Get Started",
    text: "Join now and enhance your learning journey with Quizbo today! Start achieving your comprehension goals with quizzes tailored just for you.",
  },
];

const IntroScreen = ({ onFinish }) => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const [fadeOut, setFadeOut] = useState(false);

  const handleDotClick = (index) => {
    if (index === slides.length - 1) {
      setFadeOut(true);
      setTimeout(onFinish, 1000); // Adjust as needed for fade-out duration
    } else {
      setCurrentSlide(index);
    }
  };

  return (
    <div className="intro-screen">
      <div className={`intro-container ${fadeOut ? "fade-out" : ""}`}>
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
        <div className="dots-container">
          {slides.map((_, index) => (
            <span
              key={index}
              className={`dot ${index === currentSlide ? "active" : ""}`}
              onClick={() => handleDotClick(index)}
            ></span>
          ))}
        </div>
      </div>
    </div>
  );
};

export default IntroScreen;
