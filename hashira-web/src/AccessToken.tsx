import React from "react";
import * as firebase from "./firebase";
import Header from "./Header";
import { useFetchAccessTokens } from "./hooks";

const AccessToken: React.FC<{ user: firebase.User | null | undefined }> = ({
  user,
}) => {
  const [fetchAccessTokenState, fetchAccessTokens] = useFetchAccessTokens();
  const accesstokens = fetchAccessTokenState.data;
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});

  React.useEffect(() => {
    if (user) {
      Promise.all([fetchAccessTokens(user.uid)]).catch((e) => {
        console.log("fetch error:", JSON.stringify(e));
      });
    }
  }, [user]);

  return (
    <div>
      <Header user={user} />
      {!user || !accesstokens ? <div>Loading...</div> : (
        <>
          <button
            onClick={async () => {
              await firebase.claimNewAccessToken(user.uid);
              await fetchAccessTokens(user.uid);
            }}
          >
            Generate new access token
          </button>
          <button
            disabled={((): boolean => {
              for (const v in checkedTokens) {
                if (checkedTokens[v]) {
                  return false;
                }
              }
              return true;
            })()}
            onClick={async () => {
              const accesstokens: string[] = [];
              for (let v in checkedTokens) {
                accesstokens.push(v);
              }

              await firebase.revokeAccessTokens(user.uid, accesstokens);
              await fetchAccessTokens(user.uid);
            }}
          >
            Revoke access token
          </button>
          {accesstokens.map((token: string) => {
            return (
              <li key={token}>
                <input
                  type="checkbox"
                  checked={checkedTokens[token] || false}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    setCheckedTokens({
                      ...checkedTokens,
                      [token]: e.target.checked,
                    });
                  }}
                  name={token}
                  value={token}
                />
                {token}
              </li>
            );
          })}
        </>
      )}
    </div>
  );
};

export default AccessToken;
