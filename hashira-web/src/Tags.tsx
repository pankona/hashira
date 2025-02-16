import React from "react";
import styled from "styled-components";
import * as firebase from "./firebase";
import Header from "./Header";
import { useFetchTasksAndPriorities, useUpdateTasks } from "./hooks";
import { StyledVerticalSpacer } from "./styles";
import { useNavigate } from "react-router-dom";

const StyledTagList = styled.div`
  padding: 16px;
  max-width: 800px;
  margin: 0 auto;
`;

const StyledTagItem = styled.div`
  display: grid;
  grid-template-columns: minmax(0, 1fr) 80px 40px;
  align-items: center;
  padding: 8px;
  border-radius: 4px;
  &:hover {
    background-color: #f5f5f5;
  }
`;

const StyledTagName = styled.div`
  font-weight: bold;
  font-size: 16px;
  padding-right: 16px;
  cursor: pointer;
  &:hover {
    text-decoration: underline;
  }
`;

const StyledTaskCount = styled.div`
  color: #666;
  font-size: 14px;
  text-align: left;
`;

const StyledButton = styled.button`
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  width: 32px;
  height: 32px;
  &:hover {
    background-color: #eee;
  }
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const StyledEditForm = styled.div`
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  align-items: center;
  gap: 8px;
  grid-column: 1 / -1;
`;

const StyledInput = styled.input`
  padding: 4px 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
  &:focus {
    outline: none;
    border-color: #666;
  }
`;

interface TagInfo {
  name: string;
  taskCount: number;
  tasks: {
    id: string;
    text: string;
  }[];
}

const Tags: React.FC<{ user: firebase.User | null | undefined }> = ({
  user,
}) => {
  const navigate = useNavigate();
  const [fetchTasksAndPrioritiesState, fetchTasksAndPriorities] = useFetchTasksAndPriorities();
  const [updateTasksState, updateTasks] = useUpdateTasks();
  const [editingTag, setEditingTag] = React.useState<string | null>(null);
  const [newTagName, setNewTagName] = React.useState("");

  const tasksAndPriorities = fetchTasksAndPrioritiesState.data as firebase.TasksAndPriorities | null;
  const isLoading = updateTasksState.isLoading || fetchTasksAndPrioritiesState.isLoading;

  React.useEffect(() => {
    if (user) {
      fetchTasksAndPriorities(user.uid).catch((e) => {
        console.log("fetch error:", JSON.stringify(e));
      });
    }
  }, [user]);

  const extractTags = React.useMemo(() => {
    if (!tasksAndPriorities) return [];

    const tagMap = new Map<string, TagInfo>();

    // タスクからタグを抽出
    Object.entries(tasksAndPriorities.Tasks).forEach(([taskId, task]) => {
      if (task.IsDeleted) return;

      // タグの抽出条件を修正:
      // 1. 文字列の先頭が1つ以上の # で始まる場合
      // 2. または、1つ以上の # の前にスペースがある場合
      const tags = task.Name.match(/(^|\s)#+[^\s]+/g) || [];
      tags.forEach((tag: string) => {
        // タグ名には # の数も含める
        const tagName = tag.trim();
        const existingTag = tagMap.get(tagName) || {
          name: tagName,
          taskCount: 0,
          tasks: [],
        };

        existingTag.taskCount++;
        existingTag.tasks.push({
          id: taskId,
          text: task.Name,
        });
        tagMap.set(tagName, existingTag);
      });
    });

    return Array.from(tagMap.values()).sort((a, b) => b.taskCount - a.taskCount);
  }, [tasksAndPriorities]);

  const handleTagEdit = async (oldTag: string, newTag: string) => {
    if (!user || !tasksAndPriorities) return;

    // 新しいタグ名に # を付与
    const newTagWithHash = '#' + newTag;
    if (oldTag === newTagWithHash) return;

    const tasksToUpdate: firebase.TasksObject = {};
    const tagInfo = extractTags.find((t) => t.name === oldTag);
    
    if (!tagInfo) return;

    tagInfo.tasks.forEach(({ id, text }) => {
      const task = tasksAndPriorities.Tasks[id];
      // タグをそのまま置換（# の数を含めて置換）
      tasksToUpdate[id] = {
        ...task,
        Name: text.replace(new RegExp(oldTag.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), "g"), newTagWithHash),
      };
    });

    try {
      await updateTasks(tasksToUpdate, false);
      await fetchTasksAndPriorities(user.uid);
      setEditingTag(null);
      setNewTagName("");
    } catch (e) {
      console.error("Failed to update tags:", e);
    }
  };

  // タグ名から # プレフィックスを取得
  const getHashPrefix = (tagName: string) => {
    const match = tagName.match(/^#+/);
    return match ? match[0] : '#';
  };

  // タグ名から # プレフィックスを除去
  const removeHashPrefix = (tagName: string) => {
    return tagName.replace(/^#+/, '');
  };

  const handleTagClick = (tagName: string) => {
    localStorage.setItem("hashira-filter-text", tagName);
    navigate("/");
  };

  if (!user) {
    return (
      <div>
        <Header user={user} />
        <div>Please login to view tags.</div>
      </div>
    );
  }

  if (isLoading || !tasksAndPriorities) {
    return (
      <div>
        <Header user={user} />
        <div>Loading...</div>
      </div>
    );
  }

  return (
    <div>
      <Header user={user} />
      <StyledTagList>
        <h2>Tags</h2>
        <StyledVerticalSpacer />
        {extractTags.length === 0 ? (
          <div>No tags found. Add tags to your tasks by including #tag in the task text.</div>
        ) : (
          extractTags.map((tag) => (
            <StyledTagItem key={tag.name}>
              {editingTag === tag.name ? (
                <StyledEditForm>
                  <span style={{ color: "#666" }}>{getHashPrefix(tag.name)}</span>
                  <StyledInput
                    type="text"
                    value={newTagName}
                    onChange={(e) => {
                      setNewTagName(e.target.value);
                    }}
                    autoFocus
                  />
                  <StyledButton
                    onClick={() => handleTagEdit(tag.name, newTagName)}
                    disabled={!newTagName || ('#' + newTagName === tag.name)}
                    title="Save"
                  >
                    ✓
                  </StyledButton>
                  <StyledButton
                    onClick={() => {
                      setEditingTag(null);
                      setNewTagName("");
                    }}
                    title="Cancel"
                  >
                    ✕
                  </StyledButton>
                </StyledEditForm>
              ) : (
                <>
                  <StyledTagName onClick={() => handleTagClick(tag.name)}>{tag.name}</StyledTagName>
                  <StyledTaskCount>({tag.taskCount} tasks)</StyledTaskCount>
                  <StyledButton
                    onClick={() => {
                      setEditingTag(tag.name);
                      setNewTagName(removeHashPrefix(tag.name));
                    }}
                    title="Edit tag"
                  >
                    ✎
                  </StyledButton>
                </>
              )}
            </StyledTagItem>
          ))
        )}
      </StyledTagList>
    </div>
  );
};

export default Tags; 