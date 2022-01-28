import { initializeApp } from "firebase/app";
import {
  getAuth,
  GoogleAuthProvider,
  onAuthStateChanged as onFirebaseAuthStateChanged,
  signInWithRedirect,
  signOut,
  User as FirebaseUser,
} from "firebase/auth";
import {
  addDoc,
  collection,
  deleteDoc,
  doc,
  DocumentData,
  FieldValue,
  getDoc,
  getDocs,
  getFirestore,
  query,
  QueryDocumentSnapshot,
  serverTimestamp,
  where,
} from "firebase/firestore";
import * as functions from "firebase/functions";

import { v4 as uuidv4 } from "uuid";

const firebaseConfig = {
  apiKey: "AIzaSyDMkM3qb_CUokFQDSFemLhPOqXJrR-rVbo",
  authDomain: "hashira-web.firebaseapp.com",
  projectId: "hashira-web",
  storageBucket: "hashira-web.appspot.com",
  messagingSenderId: "150558268935",
  appId: "1:150558268935:web:74eef753ffba6bb8bd54a2",
  measurementId: "G-EEZ5MJJ6XL",
};

initializeApp(firebaseConfig);

export const login = () => {
  const auth = getAuth();
  const provider = new GoogleAuthProvider();
  signInWithRedirect(auth, provider);
};

export const logout = () => {
  const auth = getAuth();
  signOut(auth);
};

export const onAuthStateChanged = (cb: (user: User | null) => void) => {
  const auth = getAuth();
  onFirebaseAuthStateChanged(auth, cb);
};

export type User = FirebaseUser;

interface accesstoken {
  uid: string;
  accesstoken: string;
  timestamp: FieldValue;
}

export const claimNewAccessToken = async (uid: string) => {
  const db = getFirestore();
  const data: accesstoken = {
    uid: uid,
    accesstoken: uuidv4(),
    timestamp: serverTimestamp(),
  };
  await addDoc(collection(db, "accesstokens"), data);
};

export const fetchAccessTokens = async (uid: string): Promise<string[]> => {
  const db = getFirestore();
  const querySnapshot = await getDocs(
    query(collection(db, "accesstokens"), where("uid", "==", uid))
  );
  const ret: accesstoken[] = [];
  querySnapshot.forEach((doc: QueryDocumentSnapshot<DocumentData>) => {
    const data = doc.data({ serverTimestamps: "estimate" });

    const token: accesstoken = {
      uid: data.uid,
      accesstoken: data.accesstoken,
      timestamp: data.timestamp,
    };
    ret.push(token);
  });

  return ret
    .sort((a: any, b: any) => {
      return a.timestamp.seconds - b.timestamp.seconds;
    })
    .map((a: accesstoken) => {
      return a.accesstoken;
    });
};

export const revokeAccessTokens = async (
  uid: string,
  accesstokens: string[]
) => {
  const db = getFirestore();

  for (const accesstoken of accesstokens) {
    const querySnapshot = await getDocs(
      query(
        collection(db, "accesstokens"),
        where("uid", "==", uid),
        where("accesstoken", "==", accesstoken)
      )
    );

    for (const doc of querySnapshot.docs) {
      await deleteDoc(doc.ref);
    }
  }
};

export const Place = ["BACKLOG", "TODO", "DOING", "DONE"] as const;

// uploadTasks
// 複数の task を受け取って、全部 BACKLOG の一番上に積む
export const uploadTasks = async (tasks: string[]) => {
  const tasksObject: {
    [key: string]: {
      ID: string;
      IsDeleted: boolean;
      Name: string;
      Place: typeof Place[number];
    };
  } = {};
  const priorities: string[] = [];

  tasks.forEach((v: string) => {
    const taskId = uuidv4();
    tasksObject[taskId] = {
      ID: taskId,
      IsDeleted: false,
      Name: v,
      Place: "BACKLOG",
    };
    priorities.push(taskId);
  });

  try {
    await functions.httpsCallable(
      functions.getFunctions(undefined, "asia-northeast1"),
      "call?method=add"
    )({
      tasks: tasksObject,
      priority: {
        BACKLOG: priorities,
      },
    });
  } catch (e) {
    // FIXME:
    // currently cloud functions doesn't return appropriate response
    // that fits httpsCallable protocol even if the function succeeded.
    console.log("error:", e);
  }
};

export interface TasksObject {
  [key: string]: {
    ID: string;
    IsDeleted: boolean;
    Name: string;
    Place: typeof Place[number];
  };
}

// updateTasks
// task の状態を変えるために用いる。変更があった task はそれぞれのレーンの一番上に積まれる。
// もっぱら、タスクの状態を変更する (横移動する) ために用いる。
export const updateTasks = async (tasksObject: TasksObject) => {
  const priorities: {
    [key: string]: string[];
  } = {};

  Place.forEach((v) => {
    priorities[v] = [];
  });

  for (const v in tasksObject) {
    priorities[tasksObject[v].Place].push(tasksObject[v].ID);
  }

  try {
    await functions.httpsCallable(
      functions.getFunctions(undefined, "asia-northeast1"),
      "call?method=add"
    )({
      tasks: tasksObject,
      priority: priorities,
    });
  } catch (e) {
    // FIXME:
    // currently cloud functions doesn't return appropriate response
    // that fits httpsCallable protocol even if the function succeeded.
    console.log("error:", e);
  }
};

// updateTasks2
// task の状態を変更するために用いる。変更があった task の場所はそのまま維持される。
// もっぱら、task の中身を編集するときに用いられる。
// FIXME: 名前がひどいので直すこと
export const updateTasks2 = async (tasksObject: TasksObject) => {
  try {
    await functions.httpsCallable(
      functions.getFunctions(undefined, "asia-northeast1"),
      "call?method=add"
    )({
      tasks: tasksObject,
    });
  } catch (e) {
    // FIXME:
    // currently cloud functions doesn't return appropriate response
    // that fits httpsCallable protocol even if the function succeeded.
    console.log("error:", e);
  }
};

export const fetchTaskAndPriorities = async (uid: string) => {
  const db = getFirestore();
  const docRef = doc(db, "tasksAndPriorities", uid);
  const docSnapshot = await getDoc(docRef);
  return docSnapshot.data();
};
