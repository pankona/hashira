import React from "react";
import styled from "styled-components";

const StyledInputForm = styled.form`
  display: flex;
`;

const TaskInput: React.VFC<{
  onSubmitTasks: (tasks: string[]) => void;
  disabled: boolean;
}> = ({ onSubmitTasks, disabled }) => {
  const [tasks, setTasks] = React.useState<string[]>([]);

  return (
    <StyledInputForm>
      <textarea
        placeholder={"Add todos"}
        value={tasks.join("\n")}
        onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
          setTasks(e.target.value.split("\n"));
        }}
      ></textarea>
      <input
        type="submit"
        value="Submit"
        autoFocus={true}
        disabled={tasks.length === 0 || disabled}
        onClick={(e: React.FormEvent<HTMLInputElement>) => {
          e.preventDefault();
          onSubmitTasks(tasks);
          setTasks([]);
        }}
      />
    </StyledInputForm>
  );
};

export default TaskInput;
