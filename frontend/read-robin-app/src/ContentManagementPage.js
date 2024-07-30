import React, { useState, useEffect } from "react";
import "./ContentManagementPage.css";
import { db, storage } from "./firebase";
import { collection, getDocs } from "firebase/firestore";
import { getDownloadURL, ref } from "firebase/storage";

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
            const data = doc.data();
            if (data.url && data.url.startsWith("gs://")) {
              const fileRef = ref(
                storage,
                data.url.replace("gs://read-robin-2e150.appspot.com/", "")
              );
              try {
                const url = await getDownloadURL(fileRef);
                data.url = url;
              } catch (error) {
                console.error("Error fetching download URL:", error);
                throw error;
              }
            }
            return {
              id: doc.id,
              ...data,
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

    console.log("Payload being sent to backend:", payload); // Log the payload

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

  return (
    <div className="content-management-page">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Your Content</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="loading-spinner"></div>}
      {!loading && contents.length > 0 && (
        <div className="content-list">
          {contents.map((content) => (
            <div key={content.id} className="content-item">
              <a href={content.url} target="_blank" rel="noopener noreferrer">
                <h3>{content.title}</h3>
              </a>
              <button
                className="generate-new-quiz"
                onClick={() =>
                  handleGenerateQuiz(
                    content.id,
                    content.title,
                    content.url,
                    content.content_text
                  )
                }
              >
                Generate Quiz
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
    </div>
  );
}

export default ContentManagementPage;
