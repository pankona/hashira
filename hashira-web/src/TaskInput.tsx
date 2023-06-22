import React, { useId } from "react";
import styled from "styled-components";
import { StyledHorizontalSpacer } from "./styles";
import { normalizeTasks } from "./task";

const StyledInputForm = styled.form`
  display: flex;
`;

const TaskInput: React.VFC<{
  onSubmitTasks: (tasks: string[]) => Promise<void>;
  disabled: boolean;
}> = ({ onSubmitTasks, disabled }) => {
  const [tasks, setTasks] = React.useState<string[]>([]);

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
    </StyledInputForm>
  );
};

export default TaskInput;
