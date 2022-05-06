import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import AccessToken from "./AccessToken";
import App from "./App";
import { useUser } from "./hooks";

const Router: React.FC = () => {
  const user = useUser();
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App user={user} />} />
        <Route path="/accesstokens" element={<AccessToken user={user} />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
