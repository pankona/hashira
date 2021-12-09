import React from "react";
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
`;

const StyledCheckbox = styled.input.attrs({ type: "checkbox" })`
  margin-right: 8px;
`;

export const TaskList: React.VFC<{
  user: firebase.User;
  place: typeof firebase.Place[number];
  tasksAndPriorities: any;
  checkedTasks: { [key: string]: boolean };
  setCheckedTasks: (a: { [key: string]: boolean }) => void;
  setTasksAndPriorities: (tp: any | undefined) => void;
}> = ({
  user,
  place,
  tasksAndPriorities,
  checkedTasks,
  setCheckedTasks,
  setTasksAndPriorities,
}) => {
  const [updatedTasks, setUpdatedTasks] = React.useState<{
    [key: string]: string;
  }>({});
  return (
    <StyledList>
      {tasksAndPriorities["Priority"][place] &&
        tasksAndPriorities["Priority"][place]
          .filter((v: any) => tasksAndPriorities["Tasks"][v])
          .map((p: string) => {
            const taskId = tasksAndPriorities["Tasks"][p].ID;
            const taskName = tasksAndPriorities["Tasks"][p].Name;
            return (
              <StyledListItem key={taskId}>
                <StyledCheckbox
                  id={taskId}
                  value={taskName}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    setCheckedTasks({
                      ...checkedTasks,
                      [e.target.id]: e.target.checked,
                    });
                  }}
                />
                <StyledListContent
                  id={taskId}
                  key={taskId}
                  value={
                    updatedTasks[taskId] !== undefined
                      ? updatedTasks[taskId]
                      : taskName
                  }
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    setUpdatedTasks({
                      ...updatedTasks,
                      [e.target.id]: e.target.value,
                    });
                  }}
                  onBlur={async (_e: React.ChangeEvent<HTMLInputElement>) => {
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
                        Name:
                          updatedTasks[v] !== "" ? updatedTasks[v] : taskName,
                        Place: task.Place,
                      };
                    }

                    await firebase.updateTasks2(tasksToUpdate);

                    // refresh tasks and priorities
                    const tp = await firebase.fetchTaskAndPriorities(user.uid);
                    setTasksAndPriorities(tp);

                    setUpdatedTasks({});
                  }}
                />
              </StyledListItem>
            );
          })}
    </StyledList>
  );
};
