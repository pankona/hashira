import React, { useRef } from "react";
import styled from "styled-components";
import * as firebase from "./firebase";
import Header from "./Header";
import { useAddTasks, useFetchTasksAndPriorities, useFilteredTasks, useUpdateTasks } from "./hooks";
import { StyledHorizontalSpacer, StyledVerticalSpacer } from "./styles";
import TaskInput from "./TaskInput";
import { TaskList } from "./TaskList";

const StyledBody = styled.div`
  padding-left: 8px;
  padding-right: 8px;
`;

const App: React.FC<{ user: firebase.User | null | undefined }> = ({
  user,
}) => {
  const [checkedTasks, setCheckedTasks] = React.useState<{
    [key: string]: boolean;
  }>({});
  const [mode, setMode] = React.useState<"move" | "select">("select");
  const [filter, setFilter] = React.useState<string>("");

  const [addTasksState, addTasks] = useAddTasks();
  const [updateTasksState, updateTasks] = useUpdateTasks();
  const [fetchTasksAndPrioritiesState, fetchTasksAndPriorities] = useFetchTasksAndPriorities();

  const tasksAndPriorities = fetchTasksAndPrioritiesState.data;
  const isLoading = addTasksState.isLoading
    || updateTasksState.isLoading
    || fetchTasksAndPrioritiesState.isLoading;

  const filteredTasks = useFilteredTasks(tasksAndPriorities ? tasksAndPriorities.Tasks : [], filter);

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
    [user],
  );

  const isMoveTaskProcessing = useRef(false);

  const onMoveTask = React.useCallback(
    (taskId: string, direction: "left" | "right") => {
      return new Promise<void>(async (resolve, reject) => {
        if (!user || isMoveTaskProcessing.current) {
          resolve();
          return;
        }
        const tasksToMove: firebase.TasksObject = {};
        const task = tasksAndPriorities["Tasks"][taskId];
        const currentIndex = firebase.Places.indexOf(task.Place);
        const nextIndex = ((): number => {
          if (direction === "left") {
            if (currentIndex === 0) {
              return firebase.Places.length - 1;
            }
            return currentIndex - 1;
          }

          if (currentIndex === firebase.Places.length - 1) {
            return 0;
          }
          return currentIndex + 1;
        })();

        tasksToMove[taskId] = {
          ID: task.ID,
          IsDeleted: false,
          Name: task.Name,
          Place: firebase.Places[nextIndex],
        };

        try {
          isMoveTaskProcessing.current = true;
          await updateTasks(tasksToMove, true);
          await fetchTasksAndPriorities(user.uid);
          resolve();
        } catch (e) {
          reject(e);
        } finally {
          isMoveTaskProcessing.current = false;
        }
      });
    },
    [user, tasksAndPriorities],
  );

  const onEditTasks = React.useCallback(
    (tasks: firebase.TasksObject) => {
      return new Promise<void>(async (resolve, reject) => {
        if (!user) {
          resolve();
          return;
        }
        try {
          await updateTasks(tasks, false);
          // refresh tasks and priorities
          await fetchTasksAndPriorities(user.uid);
          resolve();
        } catch (e) {
          reject(e);
        }
      });
    },
    [user],
  );

  const intervalMs = 1000 * 60 * 3; // 3 minutes

  React.useEffect(() => {
    firebase.ping();
    const intervalId = setInterval(() => {
      firebase.ping();
    }, intervalMs);

    return () => {
      clearInterval(intervalId);
    };
  }, []);

  React.useEffect(() => {
    if (user) {
      Promise.all([fetchTasksAndPriorities(user.uid)]).catch((e) => {
        console.log("fetch error:", JSON.stringify(e));
      });
    }
  }, [user]);

  return (
    <div>
      <Header user={user} />
      <StyledBody>
        {user !== null && (
          <TaskInput
            onSubmitTasks={onSubmitTasks}
            disabled={isLoading || !user}
            onFilterChange={setFilter}
          />
        )}
        <StyledVerticalSpacer />
        {(() => {
          if (user === null) {
            return <></>;
          }

          if (!user || !tasksAndPriorities) {
            return <div>Loading...</div>;
          }

          return (
            <>
              <div style={{ display: "flex" }}>
                <input
                  type="button"
                  value={"Mark as Done"}
                  style={{ minWidth: "128px" }}
                  disabled={isLoading
                    || ((): boolean => {
                      for (const v in checkedTasks) {
                        if (checkedTasks[v]) {
                          return false;
                        }
                      }
                      return true;
                    })()}
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

                    await updateTasks(tasksToMarkAsDone, true);
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
              {tasksAndPriorities
                ? (
                  <div
                    className="TaskAndPriorities"
                    style={{
                      display: "flex",
                      overflow: "auto",
                    }}
                  >
                    <TaskList
                      place={"BACKLOG"}
                      tasksAndPriorities={{ ...tasksAndPriorities, Tasks: filteredTasks }}
                      checkedTasks={checkedTasks}
                      setCheckedTasks={setCheckedTasks}
                      onEditTasks={onEditTasks}
                      mode={mode}
                      onMoveTask={onMoveTask}
                    />
                    <TaskList
                      place={"TODO"}
                      tasksAndPriorities={{ ...tasksAndPriorities, Tasks: filteredTasks }}
                      checkedTasks={checkedTasks}
                      setCheckedTasks={setCheckedTasks}
                      onEditTasks={onEditTasks}
                      mode={mode}
                      onMoveTask={onMoveTask}
                    />
                    <TaskList
                      place={"DOING"}
                      tasksAndPriorities={{ ...tasksAndPriorities, Tasks: filteredTasks }}
                      checkedTasks={checkedTasks}
                      setCheckedTasks={setCheckedTasks}
                      onEditTasks={onEditTasks}
                      mode={mode}
                      onMoveTask={onMoveTask}
                    />
                    <TaskList
                      place={"DONE"}
                      tasksAndPriorities={{ ...tasksAndPriorities, Tasks: filteredTasks }}
                      checkedTasks={checkedTasks}
                      setCheckedTasks={setCheckedTasks}
                      onEditTasks={onEditTasks}
                      mode={mode}
                      onMoveTask={onMoveTask}
                    />
                  </div>
                )
                : undefined}
            </>
          );
        })()}
      </StyledBody>
    </div>
  );
};

export default App;
