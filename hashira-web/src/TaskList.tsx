import React from "react";
import styled from "styled-components";
import { Place } from "./firebase";

const List = styled.div`
  min-width: 300px;
  max-width: 300px;
  padding-left: 10px;
  padding-right: 10px;
  border: solid;
`;

const ListItem = styled.li`
  list-style: none;
  white-space: nowrap;
  overflow-y: scroll;
  -ms-overflow-style: none;
  scrollbar-width: none;
  ::-webkit-scrollbar {
    display: none;
  }
`;

const Checkbox = styled.input.attrs({ type: "checkbox" })`
  margin-right: 8px;
`;

export const TaskList: React.VFC<{
  place: typeof Place[number];
  tasksAndPriorities: any;
  checkedTasks: { [key: string]: boolean };
  setCheckedTasks: (a: { [key: string]: boolean }) => void;
}> = ({ place, tasksAndPriorities, checkedTasks, setCheckedTasks }) => {
  return (
    <List>
      {tasksAndPriorities["Priority"][place]
        .filter((v: any) => tasksAndPriorities["Tasks"][v])
        .map((p: string) => {
          return (
            <ListItem key={tasksAndPriorities["Tasks"][p].ID}>
              <Checkbox
                id={tasksAndPriorities["Tasks"][p].ID}
                value={tasksAndPriorities["Tasks"][p].Name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  setCheckedTasks({
                    ...checkedTasks,
                    [e.target.id]: e.target.checked,
                  });
                }}
              />
              {tasksAndPriorities["Tasks"][p].Name}
            </ListItem>
          );
        })}
    </List>
  );
};
