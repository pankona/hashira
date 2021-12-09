import React from "react";
import * as firebase from "./firebase";
import Header from "./Header";
import { TaskList } from "./TaskList";
import TaskInput from "./TaskInput";
import styled from "styled-components";
import { StyledHorizontalSpacer, StyledVerticalSpacer } from "./styles";
import {
  useAddTasks,
  useFetchAccessTokens,
  useFetchTasksAndPriorities,
  useUpdateTasks,
  useUpdateTasks2,
  useUser,
} from "./hooks";

const StyledBody = styled.div`
  padding-left: 8px;
  padding-right: 8px;
`;

const App: React.VFC = () => {
  const [checkedTokens, setCheckedTokens] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [checkedTasks, setCheckedTasks] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [mode, setMode] = React.useState<"move" | "select">("select");

  const user = useUser();
  const [addTasksState, addTasks] = useAddTasks();
  const [updateTasksState, updateTasks] = useUpdateTasks();
  const [updateTasks2State, updateTasks2] = useUpdateTasks2();
  const [fetchAccessTokenState, fetchAccessTokens] = useFetchAccessTokens();
  const [fetchTasksAndPrioritiesState, fetchTasksAndPriorities] =
    useFetchTasksAndPriorities();

  const accesstokens = fetchAccessTokenState.data;
  const tasksAndPriorities = fetchTasksAndPrioritiesState.data;
  const isLoading =
    addTasksState.isLoading ||
    updateTasksState.isLoading ||
    updateTasks2State.isLoading ||
    fetchTasksAndPrioritiesState.isLoading;

  const onSubmitTasks = React.useCallback(
    (tasks: string[]) => {
      return new Promise<void>(async (resolve, reject) => {
        if (!user) {
          resolve();
          return;
        }

        try {
          await addTasks(tasks);
          // refresh tasks and priorities
          await fetchTasksAndPriorities(user.uid);
          resolve();
        } catch (e) {
          reject(e);
        }
      });
    },
    [user]
  );

  const onMoveTask = React.useCallback(
    (taskId: string, direction: "left" | "right") => {
      return new Promise<void>(async (resolve, reject) => {
        if (!user) {
          resolve();
          return;
        }
        const tasksToMove: firebase.TasksObject = {};
        const task = tasksAndPriorities["Tasks"][taskId];
        const currentIndex = firebase.Place.indexOf(task.Place);
        const nextIndex = ((): number => {
          if (direction === "left") {
            if (currentIndex === 0) {
              return firebase.Place.length - 1;
            }
            return currentIndex - 1;
          }

          if (currentIndex === firebase.Place.length - 1) {
            return 0;
          }
          return currentIndex + 1;
        })();

        tasksToMove[taskId] = {
          ID: task.ID,
          IsDeleted: false,
          Name: task.Name,
          Place: firebase.Place[nextIndex],
        };

        try {
          await updateTasks(tasksToMove);
          await fetchTasksAndPriorities(user.uid);
          resolve();
        } catch (e) {
          reject(e);
        }
      });
    },
    [user, tasksAndPriorities]
  );

  const onEditTasks = React.useCallback(
    (tasks: firebase.TasksObject) => {
      return new Promise<void>(async (resolve, reject) => {
        if (!user) {
          resolve();
          return;
        }
        try {
          await updateTasks2(tasks);
          // refresh tasks and priorities
          await fetchTasksAndPriorities(user.uid);
          resolve();
        } catch (e) {
          reject(e);
        }
      });
    },
    [user]
  );

  React.useEffect(() => {
    if (user) {
      Promise.all([
        fetchAccessTokens(user.uid),
        fetchTasksAndPriorities(user.uid),
      ]);
    }
  }, [user]);

  return (
    <div>
      <Header
        user={user}
        isLoading={!user || !accesstokens || !tasksAndPriorities}
      />
      <StyledBody>
        {user !== null && (
          <TaskInput
            onSubmitTasks={onSubmitTasks}
            disabled={isLoading || !user}
          />
        )}
        <StyledVerticalSpacer />
        {(() => {
          if (user === null) {
            return <></>;
          }

          if (!user || !accesstokens || !tasksAndPriorities) {
            return <div>Loading...</div>;
          }

          return (
            <>
              <div style={{ display: "flex" }}>
                <input
                  type="button"
                  value={"Mark as Done"}
                  style={{ minWidth: "128px" }}
                  disabled={
                    isLoading ||
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

                    await updateTasks(tasksToMarkAsDone);
                    // refresh tasks and priorities
                    await fetchTasksAndPriorities(user.uid);
                    setCheckedTasks({});
                  }}
                />
                <StyledHorizontalSpacer />
                <input
                  type="button"
                  value={mode === "move" ? "Finish moving" : "Move"}
                  style={{ minWidth: "128px" }}
                  disabled={isLoading}
                  onClick={() => {
                    setMode((prev) => {
                      if (prev === "move") {
                        return "select";
                      }
                      setCheckedTasks({});
                      return "move";
                    });
                  }}
                />
              </div>
              <StyledVerticalSpacer />
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
                    onEditTasks={onEditTasks}
                    mode={mode}
                    onMoveTask={onMoveTask}
                  />
                  <TaskList
                    place={"TODO"}
                    tasksAndPriorities={tasksAndPriorities}
                    checkedTasks={checkedTasks}
                    setCheckedTasks={setCheckedTasks}
                    onEditTasks={onEditTasks}
                    mode={mode}
                    onMoveTask={onMoveTask}
                  />
                  <TaskList
                    place={"DOING"}
                    tasksAndPriorities={tasksAndPriorities}
                    checkedTasks={checkedTasks}
                    setCheckedTasks={setCheckedTasks}
                    onEditTasks={onEditTasks}
                    mode={mode}
                    onMoveTask={onMoveTask}
                  />
                  <TaskList
                    place={"DONE"}
                    tasksAndPriorities={tasksAndPriorities}
                    checkedTasks={checkedTasks}
                    setCheckedTasks={setCheckedTasks}
                    onEditTasks={onEditTasks}
                    mode={mode}
                    onMoveTask={onMoveTask}
                  />
                </div>
              ) : undefined}
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
          );
        })()}
      </StyledBody>
    </div>
  );
};

export default App;
