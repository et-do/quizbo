import React from "react";
import { signInWithPopup, GoogleAuthProvider } from "firebase/auth";
import { auth } from "./firebase";

function Login() {
  const provider = new GoogleAuthProvider();

  const signIn = () => {
    signInWithPopup(auth, provider).catch((error) => {
      console.error("Error signing in: ", error);
    });
  };

  return (
    <div>
      <button onClick={signIn}>Sign in with Google</button>
    </div>
  );
}

export default Login;
