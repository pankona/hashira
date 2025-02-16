import React, { useId } from "react";
import styled from "styled-components";
import { StyledHorizontalSpacer } from "./styles";
import { normalizeTasks } from "./task";

const StyledInputForm = styled.form`
  display: flex;
`;

const TaskInput: React.FC<{
  onSubmitTasks: (tasks: string[]) => Promise<void>;
  disabled: boolean;
  onFilterChange: (filter: string) => void;
}> = ({ onSubmitTasks, disabled, onFilterChange }) => {
  const [tasks, setTasks] = React.useState<string[]>([]);
  const [filter, setFilter] = React.useState<string>("");

  return (
    <StyledInputForm>
      <textarea
        placeholder={"Add todos"}
        id={useId()}
        value={tasks.join("\n")}
        onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
          setTasks(e.target.value.split("\n"));
        }}
      >
      </textarea>
      <StyledHorizontalSpacer />
      <input
        type="submit"
        value="Submit"
        autoFocus={true}
        disabled={normalizeTasks(tasks).length === 0 || disabled}
        onClick={async (e: React.FormEvent<HTMLInputElement>) => {
          e.preventDefault();
          await onSubmitTasks(tasks);
          setTasks([]);
        }}
      />
      <StyledHorizontalSpacer />
      <input
        type="text"
        placeholder="Filter tasks"
        value={filter}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          setFilter(e.target.value);
          onFilterChange(e.target.value);
        }}
      />
    </StyledInputForm>
  );
};

export default TaskInput;
