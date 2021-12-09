import React from "react";
import * as firebase from "./firebase";
import Header from "./Header";
import { TaskList } from "./TaskList";
import TaskInput from "./TaskInput";

const App: React.VFC = () => {
  const [user, setUser] = React.useState<firebase.User | null | undefined>(
    undefined
  );
  const [accesstokens, setAccessTokens] = React.useState<string[]>([]);
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [checkedTasks, setCheckedTasks] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [tasksAndPriorities, setTasksAndPriorities] = React.useState<
    any | undefined
  >(undefined);
  const [isUploading, setIsUploading] = React.useState<boolean>(false);

  const onSubmitTasks = async (tasks: string[]) => {
    if (!user) {
      return;
    }

    const tasksToAdd = tasks;
    setIsUploading(true);

    await firebase.uploadTasks(tasksToAdd);

    // refresh tasks and priorities
    const tp = await firebase.fetchTaskAndPriorities(user.uid);
    setTasksAndPriorities(tp);

    setIsUploading(false);
  };

  React.useEffect(() => {
    firebase.onAuthStateChanged((user: firebase.User | null) => {
      if (!user) {
        setUser(null);
        return;
      }

      Promise.all([
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
      <Header user={user} />
      <TaskInput
        onSubmitTasks={onSubmitTasks}
        disabled={isUploading || !user}
      />
      {user ? (
        <>
          <input
            type="button"
            value={"Mark as Done"}
            disabled={
              isUploading ||
              ((): boolean => {
                for (const v in checkedTasks) {
                  if (checkedTasks[v]) {
                    return false;
                  }
                }
                return true;
              })()
            }
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
          <input type="button" value={"Move"} disabled={isUploading} />

          {tasksAndPriorities ? (
            <div
              className="TaskAndPriorities"
              style={{
                display: "flex",
                overflow: "auto",
              }}
            >
              <TaskList
                user={user}
                place={"BACKLOG"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
                setTasksAndPriorities={setTasksAndPriorities}
              />
              <TaskList
                user={user}
                place={"TODO"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
                setTasksAndPriorities={setTasksAndPriorities}
              />
              <TaskList
                user={user}
                place={"DOING"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
                setTasksAndPriorities={setTasksAndPriorities}
              />
              <TaskList
                user={user}
                place={"DONE"}
                tasksAndPriorities={tasksAndPriorities}
                checkedTasks={checkedTasks}
                setCheckedTasks={setCheckedTasks}
                setTasksAndPriorities={setTasksAndPriorities}
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
