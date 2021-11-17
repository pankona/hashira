import React from "react";
import styled from "styled-components";
import { Place } from "./firebase";

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

const StyledListContent = styled.div`
  display: flex;
  align-items: center;
  min-height: 24px;
  white-space: nowrap;
  overflow-y: scroll;
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
  place: typeof Place[number];
  tasksAndPriorities: any;
  checkedTasks: { [key: string]: boolean };
  setCheckedTasks: (a: { [key: string]: boolean }) => void;
}> = ({ place, tasksAndPriorities, checkedTasks, setCheckedTasks }) => {
  return (
    <StyledList>
      {tasksAndPriorities["Priority"][place]
        .filter((v: any) => tasksAndPriorities["Tasks"][v])
        .map((p: string) => {
          return (
            <StyledListItem key={tasksAndPriorities["Tasks"][p].ID}>
              <StyledCheckbox
                id={tasksAndPriorities["Tasks"][p].ID}
                value={tasksAndPriorities["Tasks"][p].Name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  setCheckedTasks({
                    ...checkedTasks,
                    [e.target.id]: e.target.checked,
                  });
                }}
              />
              <StyledListContent key={tasksAndPriorities["Tasks"][p].ID}>
                {tasksAndPriorities["Tasks"][p].Name}
              </StyledListContent>
            </StyledListItem>
          );
        })}
    </StyledList>
  );
};
