import React, { useState, useEffect } from "react";
import "./ContentManagementPage.css";
import { db } from "./firebase";
import { collection, getDocs } from "firebase/firestore";

function ContentManagementPage({
  user,
  activePersona,
  setPage,
  setContentID,
  setQuizID,
}) {
  const [contents, setContents] = useState([]);

  useEffect(() => {
    const fetchContents = async () => {
      try {
        if (!user || !activePersona || !activePersona.id) {
          throw new Error("User or active persona is not defined");
        }

        const contentsRef = collection(
          db,
          "users",
          user.uid,
          "personas",
          activePersona.id,
          "quizzes"
        );
        const contentsSnapshot = await getDocs(contentsRef);
        const contentsList = contentsSnapshot.docs.map((doc) => ({
          id: doc.id,
          ...doc.data(),
        }));
        setContents(contentsList);
      } catch (error) {
        console.error("Error fetching contents:", error);
      }
    };

    fetchContents();
  }, [user, activePersona]);

  const handleGenerateQuiz = async (contentID, contentText) => {
    const payload = {
      content_text: contentText,
      persona: {
        id: activePersona.id,
        name: activePersona.name,
        role: activePersona.role,
        language: activePersona.language,
        difficulty: activePersona.difficulty,
      },
      content_type: "TEXT",
    };

    try {
      const idToken = await user.getIdToken();
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/submit`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${idToken}`,
          },
          body: JSON.stringify(payload),
        }
      );

      if (!res.ok) {
        throw new Error(`Error generating quiz: ${res.statusText}`);
      }

      const data = await res.json();
      setContentID(data.content_id);
      setQuizID(data.quiz_id);
      setPage("quizPage");
    } catch (error) {
      console.error("Error generating quiz:", error);
    }
  };

  return (
    <div className="content-management-page">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Your Saved Contents</h2>
      <div className="content-list">
        {contents.map((content) => (
          <div key={content.id} className="content-item">
            <h3>{content.title}</h3>
            <p>{content.url}</p>
            <button
              onClick={() =>
                handleGenerateQuiz(content.contentID, content.content_text)
              }
            >
              Generate Quiz
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

export default ContentManagementPage;
