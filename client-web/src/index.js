import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import reportWebVitals from "./reportWebVitals";
import TraceProvider from "./instrumentation";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <TraceProvider>
    <React.StrictMode>
      <App />
    </React.StrictMode>
  </TraceProvider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
