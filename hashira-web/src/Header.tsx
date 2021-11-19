import React from "react";
import { revision } from "./revision";
import * as firebase from "./firebase";
import styled from "styled-components";

const StyledHeader = styled.div`
  display: flex;
  justify-content: space-between;
  padding-left: 24px;
  padding-right: 24px;
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
`;

const Header: React.VFC<{ user: firebase.User | null | undefined }> = ({
  user,
}) => {
  return (
    <StyledHeader>
      <div style={{ display: "flex", minWidth: "50%" }}>
        <div style={{ minWidth: "fit-content", marginRight: "8px" }}>
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
                <div style={{ cursor: "pointer" }} onClick={firebase.logout}>
                  Logout
                </div>
              );
          }
        })()}
      </StyledLoginLogout>
    </StyledHeader>
  );
};

export default Header;