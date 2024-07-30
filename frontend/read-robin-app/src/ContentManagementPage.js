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
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showPopup, setShowPopup] = useState(false);
  const [popupContent, setPopupContent] = useState("");

  useEffect(() => {
    const fetchContents = async () => {
      setLoading(true);
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
        setError("Error fetching contents");
      } finally {
        setLoading(false);
      }
    };

    fetchContents();
  }, [user, activePersona]);

  const handleGenerateQuiz = async (contentID, title, url, contentText) => {
    setError(null);
    setLoading(true);

    const payload = {
      content_id: contentID,
      content_text: contentText,
      title: title,
      url: url,
      persona: {
        id: activePersona.id,
        name: activePersona.name,
        role: activePersona.role,
        language: activePersona.language,
        difficulty: activePersona.difficulty,
      },
      content_type: "Text",
    };

    console.log("Payload being sent to backend:", payload);

    try {
      const idToken = await user.getIdToken();
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/regenerate-quiz`,
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
        throw new Error(`Error submitting content: ${res.statusText}`);
      }

      const data = await res.json();
      setContentID(data.content_id);
      setQuizID(data.quiz_id);

      setPage("quizPage");
      setLoading(false);
    } catch (error) {
      console.error("Error:", error);
      setError(`Error submitting content: ${error.message}`);
      setLoading(false);
    }
  };

  const handleSeeContent = (contentText) => {
    setPopupContent(contentText);
    setShowPopup(true);
  };

  const closePopup = () => {
    setShowPopup(false);
    setPopupContent("");
  };

  return (
    <div className="cmp-content-management-page">
      <button className="cmp-back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Your Content</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="cmp-loading-spinner"></div>}
      {!loading && contents.length > 0 && (
        <div className="cmp-content-list">
          {contents.map((content) => (
            <div key={content.id} className="cmp-content-item">
              <h3 style={{ color: "white" }}>{content.title}</h3>
              <button
                className="cmp-generate-new-quiz"
                onClick={() =>
                  handleGenerateQuiz(
                    content.id,
                    content.title,
                    content.url,
                    content.content_text
                  )
                }
              >
                Generate New Quiz
              </button>
              <button
                className="cmp-see-content"
                onClick={() => handleSeeContent(content.content_text)}
              >
                See Content
              </button>
            </div>
          ))}
        </div>
      )}
      {contents.length === 0 && !loading && (
        <div>
          No content available. Submit some content to generate quizzes.
        </div>
      )}
      {showPopup && (
        <div className="cmp-popup-overlay">
          <div className="cmp-popup-content">
            <button className="cmp-close-button" onClick={closePopup}>
              &times;
            </button>
            <div className="cmp-popup-scroll">
              <pre>{popupContent}</pre>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default ContentManagementPage;
