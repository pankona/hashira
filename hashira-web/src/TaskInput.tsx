import React from "react";
import * as firebase from "./firebase";
import styled from "styled-components";

const StyledInputForm = styled.div``;

const TaskInput: React.VFC<{
  onSubmitTasks: (tasks: string[]) => void;
  disabled: boolean;
}> = ({ onSubmitTasks, disabled }) => {
  const [tasks, setTasks] = React.useState<string[]>([]);

  return (
    <StyledInputForm>
      <form>
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
          onClick={async (e: React.FormEvent<HTMLInputElement>) => {
            e.preventDefault();
            onSubmitTasks(tasks);
            setTasks([]);
          }}
        />
      </form>
    </StyledInputForm>
  );
};

export default TaskInput;
