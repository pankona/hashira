import React from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import AccessToken from "./AccessToken";
import App from "./App";
import { useUser } from "./hooks";
import SearchBar from "./SearchBar"; // Import the SearchBar component

const Router: React.FC = () => {
  const user = useUser();

  React.useEffect(() => {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker
        .register("/service-worker.js")
        .then(() => {
          console.log("service worker is registered");
        })
        .catch((e) => {
          console.log("failed to register service worker: ", JSON.stringify(e));
        });
    }
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App user={user} />} />
        <Route path="/accesstokens" element={<AccessToken user={user} />} />
        <Route path="/search" element={<SearchBar onSearch={() => {}} />} /> {/* Add a route for the SearchBar component */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
