import React from "react";
import * as firebase from "./firebase";

const App: React.VFC = () => {
  const [user, setUser] = React.useState<firebase.User | null>(null);
  const [accesstokens, setAccessTokens] = React.useState<string[]>([]);
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [task, setTask] = React.useState<string>("");

  React.useEffect(() => {
    firebase.onAuthStateChanged((user: firebase.User | null) => {
      setUser(user);
      if (user) {
        (async () => {
          const ret = await firebase.fetchAccessTokens(user.uid);
          setAccessTokens(ret);
        })();
      }
    });
  }, []);

  return (
    <div>
      <button onClick={firebase.login}>Login</button>
      <button onClick={firebase.logout}>Logout</button>
      {user ? (
        <>
          <div>{`Hello, ${user?.displayName!}!`}</div>
          <form>
            <label>
              Add a todo&nbsp;
              <input
                type="text"
                name="todo"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  setTask(e.target.value);
                }}
                value={task}
                autoFocus={true}
              />
            </label>
            <input
              type="submit"
              value="Submit"
              disabled={task === ""}
              onClick={(e: React.FormEvent<HTMLInputElement>) => {
                e.preventDefault();
                firebase.UploadTasks(task);
                setTask("");
              }}
            />
          </form>
          <button
            onClick={async () => {
              await firebase.claimNewAccessToken(user.uid);

              const ret = await firebase.fetchAccessTokens(user.uid);
              setAccessTokens(ret);
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

              await firebase.RevokeAccessTokens(user.uid, accesstokens);

              const ret = await firebase.fetchAccessTokens(user.uid);
              setAccessTokens(ret);
            }}
          >
            Revoke access token
          </button>
          {accesstokens.map((token: string) => {
            return (
              <li key={token}>
                <input
                  id={token}
                  type="checkbox"
                  checked={checkedTokens[token] || false}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    setCheckedTokens({
                      ...checkedTokens,
                      [e.target.id]: e.target.checked,
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
      ) : undefined}
    </div>
  );
};

export default App;
