import React, { useState } from "react";
import "./VideoForm.css";
import { db, storage } from "./firebase";
import { doc, setDoc } from "firebase/firestore";
import { ref, uploadBytes, getDownloadURL } from "firebase/storage";

function VideoForm({ user, activePersona, setPage, setContentID, setQuizID }) {
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const exampleVideo = {
    text: "Happiness - A Very Short Story",
    url: "https://storage.googleapis.com/read-robin-examples/video/happiness_a_very_short_story.mp4",
  };

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (!user || !activePersona || !activePersona.id) {
        throw new Error("User or active persona is not defined");
      }

      // Upload file to Firebase Storage
      const storageRef = ref(storage, `${user.uid}/video/${file.name}`);
      await uploadBytes(storageRef, file);
      const fileURL = await getDownloadURL(storageRef);

      // Construct the GCS URI
      const bucketName = "read-robin-2e150.appspot.com";
      const filePath = storageRef.fullPath;
      const gcsURI = `gs://${bucketName}/${filePath}`;

      // Prepare payload for submission
      const payload = {
        url: gcsURI,
        persona: {
          id: activePersona.id,
          name: activePersona.name,
          role: activePersona.role,
          language: activePersona.language,
          difficulty: activePersona.difficulty,
        },
        content_type: "Video",
      };

      // Log payload to console
      console.log("Payload being sent to backend:", payload);

      const idToken = await user.getIdToken();
      const res = await fetch(
        `https://read-robin-dev-6yudia4zva-nn.a.run.app/submit`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${idToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        }
      );

      if (!res.ok) {
        throw new Error(`Error submitting video: ${res.statusText}`);
      }

      const data = await res.json();
      setContentID(data.content_id);
      setQuizID(data.quiz_id);

      const quizRef = doc(
        db,
        "users",
        user.uid,
        "personas",
        activePersona.id,
        "quizzes",
        data.content_id
      );
      await setDoc(quizRef, {
        contentID: data.content_id,
        video_url: fileURL,
        title: data.title,
      });

      setPage("quizPage");
      setLoading(false);
    } catch (error) {
      console.error("Error:", error);
      setError(`Error submitting video: ${error.message}`);
      setLoading(false);
    }
  };

  return (
    <div className="video-form">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>Video Quiz</h2>
      <div className="example-video">
        <div className="example-card">
          <h3>Try this example:</h3>
          <ul>
            <li>
              <a
                href={exampleVideo.url}
                target="_blank"
                rel="noopener noreferrer"
                className="example-link"
              >
                {exampleVideo.text}
              </a>
              {" - "}
              <a href={exampleVideo.url} download className="download-button">
                Download
              </a>
            </li>
          </ul>
        </div>
      </div>
      <form onSubmit={handleSubmit}>
        <input
          type="file"
          accept="video/*"
          onChange={handleFileChange}
          required
        />
        <button type="submit" disabled={!file}>
          Submit
        </button>
      </form>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {loading && <div className="loading-spinner"></div>}
    </div>
  );
}

export default VideoForm;
