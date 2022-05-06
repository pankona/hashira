import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import AccessToken from "./AccessToken";
import App from "./App";

const Router: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/accesstokens" element={<AccessToken />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
