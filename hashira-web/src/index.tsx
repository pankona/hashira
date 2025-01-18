import React from "react";
import { createRoot } from "react-dom/client";
import { createGlobalStyle } from "styled-components";
import { normalize } from "styled-normalize";
import Router from "./Router";
import { assertIsDefined } from "./types";

const GlobalStyle = createGlobalStyle`
  ${normalize}
`;

const rootElement = document.getElementById("app");
assertIsDefined(rootElement);
createRoot(rootElement).render(
  <React.StrictMode>
    <GlobalStyle />
    <Router />
  </React.StrictMode>,
);
