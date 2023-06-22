import React from "react";
import { createRoot } from "react-dom/client";
import { Normalize } from "styled-normalize";
import Router from "./Router";
import { assertIsDefined } from "./types";

const rootElement = document.getElementById("app");
assertIsDefined(rootElement);
createRoot(rootElement).render(
  <React.StrictMode>
    <Normalize />
    <Router />
  </React.StrictMode>,
);
