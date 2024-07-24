import React from "react";
import "./SelectionPage.css";

function SelectionPage({ setPage }) {
  return (
    <div className="selection-page">
      <div className="selection-options">
        <h1>What do you want a quiz generated for?</h1>
        <button className="selection-button" onClick={() => setPage("urlForm")}>
          Webpage
        </button>
        <button className="selection-button" onClick={() => setPage("pdfForm")}>
          PDF
        </button>
        <button
          className="selection-button"
          onClick={() => setPage("audioForm")}
        >
          Audio
        </button>
      </div>
    </div>
  );
}

export default SelectionPage;
