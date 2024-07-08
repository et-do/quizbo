import React from "react";

function Login({ signIn }) {
  return (
    <div>
      <button onClick={signIn}>Sign in with Google</button>
    </div>
  );
}

export default Login;
