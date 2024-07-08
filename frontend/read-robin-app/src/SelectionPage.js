import React from "react";
import "./App.css";

function SelectionPage({ setPage }) {
  return (
    <div className="selection-options">
      <h1>What do you want a quiz generated for?</h1>
      <button className="selection-button" onClick={() => setPage("quizForm")}>
        Webpage
      </button>
      <button className="selection-button" disabled>
        PDF (Coming Soon)
      </button>
      <button className="selection-button" disabled>
        Audio (Coming Soon)
      </button>
    </div>
  );
}

export default SelectionPage;
