import React from "react";
import * as firebase from "./firebase";

const App: React.VFC = () => {
  const [user, setUser] = React.useState<firebase.User | null>(null);
  const [accesstokens, setAccessTokens] = React.useState<string[]>([]);

  React.useEffect(() => {
    console.log("effect");
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
          <button
            onClick={() => {
              firebase.claimNewAccessToken(user.uid);
              (async () => {
                const ret = await firebase.fetchAccessTokens(user.uid);
                setAccessTokens(ret);
              })();
            }}
          >
            Generate new access token
          </button>
          {accesstokens.map((token: string) => {
            return <li key={token}>{token}</li>;
          })}
        </>
      ) : undefined}
    </div>
  );
};

export default App;
