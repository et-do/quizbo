import React from "react";
import "./Login.css";

function Login({ signIn }) {
  return (
    <div className="login-page">
      <button className="login-button" onClick={signIn}>
        Sign in with Google
      </button>
    </div>
  );
}

export default Login;
