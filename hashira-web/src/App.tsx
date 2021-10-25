import React from "react";
import * as firebase from "./firebase";
import styled from "styled-components";

const TaskListItem = styled.li`
  white-space: nowrap;
  overflow-y: scroll;
  -ms-overflow-style: none;
  scrollbar-width: none;
  ::-webkit-scrollbar {
    display: none;
  }
`;

const TaskList = styled.div`
  width: 300px;
  padding-left: 10px;
  padding-right: 10px;
  border: solid;
`;

const App: React.VFC = () => {
  const [user, setUser] = React.useState<firebase.User | null>(null);
  const [accesstokens, setAccessTokens] = React.useState<string[]>([]);
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [task, setTask] = React.useState<string>("");
  const [tasksAndPriorities, setTasksAndPriorities] = React.useState<
    any | undefined
  >(undefined);

  React.useEffect(() => {
    firebase.onAuthStateChanged((user: firebase.User | null) => {
      setUser(user);
      if (user) {
        (async () => {
          const ret = await firebase.fetchAccessTokens(user.uid);
          setAccessTokens(ret);
          const tasksAndPriorities = await firebase.fetchTaskAndPriorities(
            user.uid
          );
          setTasksAndPriorities(tasksAndPriorities);
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
              onClick={async (e: React.FormEvent<HTMLInputElement>) => {
                e.preventDefault();
                const taskToAdd = task;
                setTask("");
                await firebase.uploadTasks(taskToAdd);

                // refresh tasks and priorities
                const tasksAndPriorities =
                  await firebase.fetchTaskAndPriorities(user.uid);
                setTasksAndPriorities(tasksAndPriorities);
              }}
            />
          </form>
          {tasksAndPriorities ? (
            <div
              className="TaskAndPriorities"
              style={{
                display: "flex",
                overflow: "auto",
                width: "min-content",
              }}
            >
              <TaskList>
                {tasksAndPriorities["Priority"]["BACKLOG"].map((p: string) => {
                  return (
                    <TaskListItem key={p}>
                      {tasksAndPriorities["Tasks"][p].Name}
                    </TaskListItem>
                  );
                })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["TODO"].map((p: string) => {
                  return (
                    <TaskListItem key={p}>
                      {tasksAndPriorities["Tasks"][p].Name}
                    </TaskListItem>
                  );
                })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["DOING"].map((p: string) => {
                  return (
                    <TaskListItem key={p}>
                      {tasksAndPriorities["Tasks"][p].Name}
                    </TaskListItem>
                  );
                })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["DONE"].map((p: string) => {
                  return (
                    <TaskListItem key={p}>
                      {tasksAndPriorities["Tasks"][p].Name}
                    </TaskListItem>
                  );
                })}
              </TaskList>
            </div>
          ) : undefined}
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

              await firebase.revokeAccessTokens(user.uid, accesstokens);

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
