import React from "react";
import * as firebase from "./firebase";
import { revision } from "./revision";
import { TaskList } from "./TaskList";

const App: React.VFC = () => {
  const [user, setUser] = React.useState<firebase.User | null>(null);
  const [accesstokens, setAccessTokens] = React.useState<string[]>([]);
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [tasks, setTasks] = React.useState<string[]>([]);
  const [checkedTasks, setCheckedTasks] = React.useState<{
    [key: string]: boolean;
  }>({});
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
      <form>
        <label>
          Add todos&nbsp;
          <textarea
            value={tasks.join("\n")}
            onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
              setTasks(e.target.value.split("\n"));
            }}
          ></textarea>
        </label>
        <input
          type="submit"
          value="Submit"
          autoFocus={true}
          disabled={tasks.length === 0 || isUploading || !user}
          onClick={async (e: React.FormEvent<HTMLInputElement>) => {
            e.preventDefault();
            if (!user) {
              return;
            }

            const tasksToAdd = tasks;
            setTasks([]);
            setIsUploading(true);

            await firebase.uploadTasks(tasksToAdd);

            // refresh tasks and priorities
            const tp = await firebase.fetchTaskAndPriorities(user.uid);
            setTasksAndPriorities(tp);

            setIsUploading(false);
          }}
        />
      </form>
      <input
        type="button"
        value={"Mark as Done"}
        onClick={async (e: React.FormEvent<HTMLInputElement>) => {
          e.preventDefault();
          if (!user) {
            return;
          }

          const tasksToMarkAsDone: firebase.TasksObject = {};
          for (const v in checkedTasks) {
            if (checkedTasks[v]) {
              const task = tasksAndPriorities["Tasks"][v];
              tasksToMarkAsDone[v] = {
                ID: task.ID,
                IsDeleted: task.Place === "DONE",
                Name: task.Name,
                Place: "DONE",
              };
            }
          }

          setIsUploading(true);

          await firebase.updateTasks(tasksToMarkAsDone);

          // refresh tasks and priorities
          const tp = await firebase.fetchTaskAndPriorities(user.uid);
          setTasksAndPriorities(tp);
          setCheckedTasks({});

          setIsUploading(false);
        }}
      />

      {user ? (
        <>
          {tasksAndPriorities ? (
            <div
              className="TaskAndPriorities"
              style={{
                display: "flex",
                overflow: "auto",
              }}
            >
              <TaskList
                place={"BACKLOG"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
              />
              <TaskList
                place={"TODO"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
              />
              <TaskList
                place={"DOING"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
              />
              <TaskList
                place={"DONE"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
              />
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
