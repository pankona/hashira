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
