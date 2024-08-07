import React, { useState } from "react";
import "./PdfForm.css";
import { db, storage } from "./firebase";
import { doc, setDoc } from "firebase/firestore";
import { ref, uploadBytes, getDownloadURL } from "firebase/storage";

function PdfForm({ user, activePersona, setPage, setContentID, setQuizID }) {
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const examplePdf = {
    text: "Chemistry Chapter Page",
    url: "https://storage.googleapis.com/read-robin-examples/pdfs/chemistry_chapter_page.pdf",
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
      const storageRef = ref(storage, `${user.uid}/pdfs/${file.name}`);
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
        content_type: "PDF",
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
        throw new Error(`Error submitting PDF: ${res.statusText}`);
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
        url: data.url, // Use the normalized URL from the response
        title: data.title,
        content_text: data.content_text,
        content_type: "PDF", // Add this line
      });

      setPage("quizPage");
      setLoading(false);
    } catch (error) {
      console.error("Error:", error);
      setError(`Error submitting PDF: ${error.message}`);
      setLoading(false);
    }
  };

  return (
    <div className="pdf-form">
      <button className="back-button" onClick={() => setPage("selection")}>
        Back
      </button>
      <h2>PDF Quiz</h2>
      <div className="example-pdf">
        <div className="example-card">
          <h3>Try this example:</h3>
          <ul>
            <li>
              <a
                href={examplePdf.url}
                target="_blank"
                rel="noopener noreferrer"
                className="example-link"
              >
                {examplePdf.text}
              </a>
              {" - "}
              <a href={examplePdf.url} download className="download-button">
                Download
              </a>
            </li>
          </ul>
        </div>
      </div>
      <form onSubmit={handleSubmit}>
        <input
          type="file"
          accept="application/pdf"
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

export default PdfForm;
