import React from "react";
import * as firebase from "./firebase";
import styled from "styled-components";
import { revision } from "./revision";

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
  const [isUploading, setIsUploading] = React.useState<boolean>(false);

  React.useEffect(() => {
    firebase.onAuthStateChanged((user: firebase.User | null) => {
      if (!user) {
        setUser(null);
        return;
      }

      Promise.all<string[], any>([
        firebase.fetchAccessTokens(user.uid),
        firebase.fetchTaskAndPriorities(user.uid),
      ]).then(([accesstokens, tasksAndPriorities]) => {
        setUser(user);
        setAccessTokens(accesstokens);
        setTasksAndPriorities(tasksAndPriorities);
      });
    });
  }, []);

  return (
    <div>
      <div>hashira web {revision()}</div>
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
              disabled={task === "" || isUploading}
              onClick={async (e: React.FormEvent<HTMLInputElement>) => {
                e.preventDefault();
                const taskToAdd = task;
                setTask("");
                setIsUploading(true);
                await firebase.uploadTasks(taskToAdd);

                // refresh tasks and priorities
                const tasksAndPriorities =
                  await firebase.fetchTaskAndPriorities(user.uid);
                setTasksAndPriorities(tasksAndPriorities);

                setIsUploading(false);
              }}
            />
          </form>
          {tasksAndPriorities ? (
            <div
              className="TaskAndPriorities"
              style={{
                display: "flex",
                overflow: "auto",
              }}
            >
              <TaskList>
                {tasksAndPriorities["Priority"]["BACKLOG"]
                  .filter((v: any) => tasksAndPriorities["Tasks"][v])
                  .map((p: string) => {
                    return (
                      <TaskListItem key={p}>
                        {tasksAndPriorities["Tasks"][p].Name}
                      </TaskListItem>
                    );
                  })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["TODO"]
                  .filter((v: any) => tasksAndPriorities["Tasks"][v])
                  .map((p: string) => {
                    return (
                      <TaskListItem key={p}>
                        {tasksAndPriorities["Tasks"][p].Name}
                      </TaskListItem>
                    );
                  })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["DOING"]
                  .filter((v: any) => tasksAndPriorities["Tasks"][v])
                  .map((p: string) => {
                    return (
                      <TaskListItem key={p}>
                        {tasksAndPriorities["Tasks"][p].Name}
                      </TaskListItem>
                    );
                  })}
              </TaskList>
              <TaskList>
                {tasksAndPriorities["Priority"]["DONE"]
                  .filter((v: any) => tasksAndPriorities["Tasks"][v])
                  .map((p: string) => {
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
      ) : (
        <div>Loading...</div>
      )}
    </div>
  );
};

export default App;
