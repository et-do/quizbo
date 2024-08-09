import React from "react";
import "./App.css";

function About() {
  return (
    <div className="about-page">
      <h1>About Quizbo</h1>
      <p>
        Quizbo is your AI companion for smarter comprehension, leveraging the
        power of Gemini, Google's latest Large Language Model, to create
        tailored quizzes from a vast array of internet content.
      </p>

      <h2>Why LLMs for Quiz Generation?</h2>
      <p>
        Creating quizzes for all the internet's content is an impossible task to
        manage. However, with LLMs, not only can we create quizzes based on
        multimodal content, but we can also have those quizzes tailored
        specifically to your goals! With Quizbo and Gemini, answers can come in
        a multitude of forms as long as they capture the essence of the answer,
        moving away from the clunky days of requiring exact formats.
      </p>

      <h2>How It Works</h2>
      <h3>Content Extraction</h3>
      <p>
        Our content extractor agent behind Quizbo takes your uploaded content
        and simplifies it into plain text.
      </p>

      <h3>Quiz Generation</h3>
      <p>
        The simplified text is passed to our quiz generator agent. This agent
        creates a quiz comprising questions, correct answers, and references to
        the text to show how it derived the answers.
      </p>

      <h3>Tailoring Questions</h3>
      <p>
        Quizbo uses user-defined personas to tailor quizzes to each user's
        characteristics, difficulty, and language preferences. By defining a
        persona, users can specify their role (e.g., student, researcher),
        preferred difficulty level, and desired language. This information
        allows Quizbo to generate quizzes that are better suited to the user's
        needs and learning goals.
      </p>

      <h3>Answer Review</h3>
      <p>
        When you respond to a question, your answer, along with the content
        text, is sent to Quizbo's reviewer agent. This agent is responsible for
        grading your answer as pass or fail and providing helpful insights.
      </p>

      <h2>The Power of LLMs</h2>
      <p>
        Gemini powers every core component of Quizbo, demonstrating its immense
        capability and adaptability. By dynamically creating quizzes tailored to
        individual learning needs, from any type of content, Gemini transforms
        learning into a more interactive and impactful experience.
      </p>

      <h2>GitHub Repository</h2>
      <p>
        You can find the source code and contribute to Quizbo on our GitHub
        repository:
        <a
          href="https://github.com/et-do/quizbo"
          target="_blank"
          rel="noopener noreferrer"
        >
          https://github.com/et-do/quizbo
        </a>
      </p>
    </div>
  );
}

export default About;
