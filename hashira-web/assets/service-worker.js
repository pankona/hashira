self.addEventListener("install", (e) => {
  console.log("install service worker:", JSON.stringify(e));
});

self.addEventListener("activate", (e) => {
  console.log("activate service worker:", JSON.stringify(e));
});

self.addEventListener("fetch", () => {});
