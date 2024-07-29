import React, { useState, useEffect } from "react";
import "./ContentManagementPage.css";
import { db } from "./firebase";
import { collection, getDocs, doc } from "firebase/firestore";

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

  const handleGenerateQuiz = (contentID, quizID) => {
    setContentID(contentID);
    setQuizID(quizID);
    setPage("quizPage");
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
                handleGenerateQuiz(content.contentID, content.quizID)
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
