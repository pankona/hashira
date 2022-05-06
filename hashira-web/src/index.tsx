import React from "react";
import ReactDOM from "react-dom";
import { Normalize } from "styled-normalize";
import Router from "./Router";

const Root = () => (
  <React.Fragment>
    <Normalize />
    <Router />
  </React.Fragment>
);

ReactDOM.render(<Root />, document.getElementById("app"));

if ("serviceWorker" in navigator) {
  navigator.serviceWorker.register("/serviceworker.js").then(function () {
    console.log("Service Worker is registered!!");
  });
}
