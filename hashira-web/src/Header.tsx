import React from "react";
import { useNavigate } from "react-router-dom";
import styled from "styled-components";
import * as firebase from "./firebase";
import { revision } from "./revision";

const StyledHeader = styled.div`
  display: flex;
  justify-content: space-between;
  padding-left: 8px;
  padding-right: 8px;
  padding-top: 8px;
  padding-bottom: 8px;
`;

const StyledRevision = styled.div`
  width: 100%;
  white-space: nowrap;
  overflow-x: scroll;
  -ms-overflow-style: none;
  scrollbar-width: none;
  ::-webkit-scrollbar {
    display: none;
  }
`;

const StyledLoginLogout = styled.div`
  display: flex;
  justify-content: end;
  gap: 16px;
`;

const Header: React.FC<{
  user: firebase.User | null | undefined;
}> = ({ user }) => {
  const navigate = useNavigate();

  return (
    <StyledHeader>
      <div style={{ display: "flex", minWidth: "50%" }}>
        <div
          style={{
            minWidth: "fit-content",
            marginRight: "8px",
            cursor: "pointer",
          }}
          onClick={() => navigate("/")}
        >
          hashira web
        </div>
        <StyledRevision>{revision()}</StyledRevision>
      </div>
      <StyledLoginLogout style={{ minWidth: "50%" }}>
        {(() => {
          switch (user) {
            case undefined:
              return <div></div>;
            case null:
              return (
                <div style={{ cursor: "pointer" }} onClick={firebase.login}>
                  Login
                </div>
              );
            default:
              return (
                <>
                  <div
                    onClick={() => navigate("/tags")}
                    style={{ cursor: "pointer" }}
                  >
                    Tags
                  </div>
                  <div
                    onClick={() => navigate("/accesstokens")}
                    style={{ cursor: "pointer" }}
                  >
                    Access tokens
                  </div>
                  <div style={{ cursor: "pointer" }} onClick={firebase.logout}>
                    Logout
                  </div>
                </>
              );
          }
        })()}
      </StyledLoginLogout>
    </StyledHeader>
  );
};

export default Header;
