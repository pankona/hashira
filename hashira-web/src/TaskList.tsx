import React, { useId } from "react";
import styled from "styled-components";
import * as firebase from "./firebase";

const StyledList = styled.div`
  min-width: 300px;
  max-width: 300px;
  max-height: 80vh;
  overflow-y: auto;
  padding-left: 10px;
  padding-right: 10px;
  border: solid;
`;

const StyledListItem = styled.div`
  display: flex;
  align-items: center;
  position: relative;
  min-height: 24px;
`;

const StyledCheckbox = styled.input.attrs({ type: "checkbox" })`
  position: absolute;
`;

const StyledListContent = styled.input.attrs({ type: "text" })`
  display: flex;
  align-items: center;
  min-height: 24px;
  width: 100%;
  white-space: nowrap;
  overflow-y: scroll;
  border: none;
  -ms-overflow-style: none;
  scrollbar-width: none;
  ::-webkit-scrollbar {
    display: none;
  }
  margin-left: 24px;
  z-index: 1000;
`;

const StyledArrow = styled.div`
  position: absolute;
  min-width: 24px;
  min-height: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  font-size: 12px;
`;

export const TaskList: React.VFC<{
  place: typeof firebase.Place[number];
  tasksAndPriorities: firebase.TasksAndPriorities;
  checkedTasks: { [key: string]: boolean };
  setCheckedTasks: (a: { [key: string]: boolean }) => void;
  onEditTasks: (tasks: firebase.TasksObject) => Promise<void>;
  mode: "move" | "select";
  onMoveTask: (taskId: string, direction: "left" | "right") => Promise<void>;
}> = ({
  place,
  tasksAndPriorities,
  checkedTasks,
  setCheckedTasks,
  onEditTasks,
  mode,
  onMoveTask,
}) => {
  const [updatedTasks, setUpdatedTasks] = React.useState<{
    [key: string]: string;
  }>({});

  // onEditCompleted は、どれかひとつの task の変更が終わると呼び出される (onBlur)。
  //
  // validation の後、string を firebase.TasksObject の形に変換して onEditTasks を呼び出す。
  //
  // 典型的な使い方をすると一度に更新される task はたかだかひとつであるが、onEditTasks 内で
  // 通信エラーが起こるなどすると、変更が DB に反映されていない task が取り残される可能性がありそう。
  // そのため、変更候補のものはすべていっぺんに更新するような処理にしている。
  const onEditCompleted = async () => {
    const tasksToUpdate: firebase.TasksObject = {};
    for (const v in updatedTasks) {
      if (updatedTasks[v] === "") {
        delete updatedTasks[v];
        setUpdatedTasks({
          ...updatedTasks,
        });
        return;
      }

      const task = tasksAndPriorities["Tasks"][v];
      tasksToUpdate[v] = {
        ID: task.ID,
        IsDeleted: false,
        Name: updatedTasks[v],
        Place: task.Place,
      };
    }

    await onEditTasks(tasksToUpdate);
    setUpdatedTasks({});
  };

  const noItem = [
    <StyledListItem key="noitem" style={{ opacity: "0.3" }}>
      No item
    </StyledListItem>,
  ];

  const convertTasksAndPrioritiesToJSXElement = (
    tasksAndPriorities: firebase.TasksAndPriorities,
  ): JSX.Element[] => {
    if (!tasksAndPriorities["Priority"][place]) {
      return noItem;
    }
    const filteredItems = tasksAndPriorities["Priority"][place].filter(
      (v: string) => tasksAndPriorities["Tasks"][v],
    );
    if (filteredItems.length === 0) {
      return noItem;
    }

    return filteredItems.map((p: string) => {
      const taskId = tasksAndPriorities["Tasks"][p].ID;
      const taskName = tasksAndPriorities["Tasks"][p].Name;
      return (
        <StyledListItem key={taskId} data-tasktoken={taskId}>
          <>
            <StyledListContent
              style={{
                marginRight: mode === "select" ? "0px" : "24px",
              }}
              id={useId()}
              key={taskId}
              value={updatedTasks[taskId] !== undefined
                ? updatedTasks[taskId]
                : taskName}
              onChange={(e) => {
                setUpdatedTasks({
                  ...updatedTasks,
                  [taskId]: e.target.value,
                });
              }}
              onBlur={() => onEditCompleted()}
            />

            {mode === "select"
              ? (
                <StyledCheckbox
                  id={useId()}
                  value={taskName}
                  onChange={(e) => {
                    setCheckedTasks({
                      ...checkedTasks,
                      [taskId]: e.target.checked,
                    });
                  }}
                />
              )
              : (
                <StyledArrow>
                  <div
                    style={{ cursor: "pointer" }}
                    onClick={(e) => onMoveTask(taskId, "left")}
                  >
                    👈
                  </div>
                  <div
                    style={{ cursor: "pointer" }}
                    onClick={(e) => onMoveTask(taskId, "right")}
                  >
                    👉
                  </div>
                </StyledArrow>
              )}
          </>
        </StyledListItem>
      );
    });
  };

  return (
    <StyledList>
      <div
        className="PlaceTitle"
        style={{ minHeight: "24px", display: "flex", alignItems: "center" }}
      >
        {place}
      </div>
      {convertTasksAndPrioritiesToJSXElement(tasksAndPriorities)}
    </StyledList>
  );
};
