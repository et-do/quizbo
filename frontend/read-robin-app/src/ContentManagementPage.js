import React, { useState, useEffect } from "react";
import "./ContentManagementPage.css";
import { db } from "./firebase";
import {
  collection,
  getDocs,
  doc,
  query,
  orderBy,
  limit,
} from "firebase/firestore";

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
  const [popupTitle, setPopupTitle] = useState("");

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
        const contentsList = await Promise.all(
          contentsSnapshot.docs.map(async (doc) => {
            const attemptsRef = collection(doc.ref, "attempts");
            const attemptsSnapshot = await getDocs(attemptsRef);
            const attempts = attemptsSnapshot.docs.map((attemptDoc) => ({
              id: attemptDoc.id,
              ...attemptDoc.data(),
            }));
            const mostRecentAttempt =
              attemptsSnapshot.docs.length > 0
                ? attemptsSnapshot.docs.sort(
                    (a, b) =>
                      b.data().createdAt.seconds - a.data().createdAt.seconds
                  )[0]
                : null;
            return {
              id: doc.id,
              ...doc.data(),
              attempts: attempts.length,
              mostRecentScore: mostRecentAttempt
                ? mostRecentAttempt.data().score
                : null,
            };
          })
        );
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

  const handleSeeContent = (contentText, contentTitle) => {
    setPopupContent(contentText);
    setPopupTitle(contentTitle);
    setShowPopup(true);
  };

  const closePopup = () => {
    setShowPopup(false);
    setPopupContent("");
    setPopupTitle("");
  };

  // Group contents by content_type
  const groupContentsByContentType = (contents) => {
    return contents.reduce((acc, content) => {
      const { content_type } = content;
      if (!acc[content_type]) {
        acc[content_type] = [];
      }
      acc[content_type].push(content);
      return acc;
    }, {});
  };

  const groupedContents = groupContentsByContentType(contents);

  return (
    <div className="cmp-content-management-page">
      <h2>Your Content</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="cmp-loading-spinner"></div>}
      {!loading && Object.keys(groupedContents).length > 0 && (
        <div className="cmp-content-list">
          {Object.keys(groupedContents).map((contentType) => (
            <div key={contentType}>
              <h3 className="cmp-content-type-title">{contentType}</h3>
              {groupedContents[contentType].map((content) => (
                <div key={content.id} className="cmp-content-item">
                  <h3 style={{ color: "white" }}>{content.title}</h3>
                  <p>Attempts: {content.attempts}</p>
                  <p>
                    Most Recent Score:{" "}
                    {content.mostRecentScore !== "N/A"
                      ? `${content.mostRecentScore}%`
                      : "0%"}
                  </p>

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
                    onClick={() =>
                      handleSeeContent(content.content_text, content.title)
                    }
                  >
                    See Content
                  </button>
                </div>
              ))}
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
            <h3 className="cmp-popup-title">{popupTitle}</h3>
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
