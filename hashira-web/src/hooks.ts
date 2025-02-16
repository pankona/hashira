import React from "react";
import * as firebase from "./firebase";
import { tasksAndPrioritiesInitialValue } from "./firebase";
import { normalizeTasks } from "./task";

type APIState<T> = {
  isLoading: boolean;
  error: string | null;
  data: T | null;
};

export const useAddTasks = (): [
  APIState<void>,
  (tasksToAdd: string[]) => Promise<void>,
] => {
  const [state, setState] = React.useState<APIState<void>>({
    isLoading: false,
    error: null,
    data: null,
  });

  return [
    state,
    React.useCallback((tasksToAdd: string[]): Promise<void> => {
      return new Promise<void>((resolve, reject) => {
        setState({
          isLoading: true,
          error: null,
          data: null,
        });

        firebase
          .addTasks(normalizeTasks(tasksToAdd))
          .then(() => {
            setState({
              isLoading: false,
              error: null,
              data: null,
            });
            resolve();
          })
          .catch((e) => reject(e));
      });
    }, []),
  ];
};

export const useUpdateTasks = (): [
  APIState<void>,
  (tasks: firebase.TasksObject, updatePosition: boolean) => Promise<void>,
] => {
  const [state, setState] = React.useState<APIState<void>>({
    isLoading: false,
    error: null,
    data: null,
  });

  return [
    state,
    React.useCallback((tasks: firebase.TasksObject, updatePosition: boolean): Promise<void> => {
      return new Promise<void>((resolve, reject) => {
        setState({
          isLoading: true,
          error: null,
          data: null,
        });

        firebase
          .updateTasks(tasks, updatePosition)
          .then((result) => {
            setState({
              isLoading: false,
              error: null,
              data: result,
            });
            resolve(result);
          })
          .catch((e) => reject(e));
      });
    }, []),
  ];
};

export const useFetchTasksAndPriorities = (): [
  APIState<any>,
  (userId: string) => Promise<any>,
] => {
  const [state, setState] = React.useState<APIState<any>>({
    isLoading: false,
    error: null,
    data: null,
  });

  return [
    state,
    React.useCallback((userId: string): Promise<any> => {
      return new Promise<any>((resolve, reject) => {
        setState((prev) => {
          return {
            ...prev,
            isLoading: true,
            error: null,
          };
        });

        firebase
          .fetchTaskAndPriorities(userId)
          .then((result) => {
            if (!result) {
              // in case for empty result
              setState({
                isLoading: false,
                error: null,
                data: tasksAndPrioritiesInitialValue,
              });
              resolve(tasksAndPrioritiesInitialValue);
              return;
            }

            setState({
              isLoading: false,
              error: null,
              data: result,
            });
            resolve(result);
          })
          .catch((e) => {
            setState({
              isLoading: false,
              error: JSON.stringify(e),
              data: null,
            });
            reject(e);
          });
      });
    }, []),
  ];
};

export const useFetchAccessTokens = (): [
  APIState<string[]>,
  (userId: string) => Promise<string[]>,
] => {
  const [state, setState] = React.useState<APIState<string[]>>({
    isLoading: false,
    error: null,
    data: null,
  });

  return [
    state,
    React.useCallback((userId: string): Promise<string[]> => {
      return new Promise<string[]>((resolve, reject) => {
        setState({
          isLoading: true,
          error: null,
          data: null,
        });

        firebase
          .fetchAccessTokens(userId)
          .then((result) => {
            setState({
              isLoading: false,
              error: null,
              data: result,
            });
            resolve(result);
          })
          .catch((e) => reject(e));
      });
    }, []),
  ];
};

export const useUser = () => {
  const [state, setState] = React.useState<any | null | undefined>(undefined);

  React.useEffect(() => {
    const cachedUser = localStorage.getItem("user");
    if (cachedUser) {
      setState(JSON.parse(cachedUser));
    }

    firebase.onAuthStateChanged((user: firebase.User | null) => {
      if (!user) {
        setState(null);
        localStorage.removeItem("user");
        return;
      }
      setState(user);
      localStorage.setItem("user", JSON.stringify(user));
    });
  }, []);

  return state;
};

export const useFilteredTasks = (tasks: any, filter: string): any => {
  return React.useMemo(() => {
    if (!filter) return tasks;

    const filterWords = filter.split(" ").map(word => word.toLowerCase());

    return tasks.filter((task: any) => {
      return filterWords.every(word => task.name.toLowerCase().includes(word));
    });
  }, [tasks, filter]);
};
