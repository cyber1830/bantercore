let token = "";
const list = document.querySelector("#discussions");
const count = document.querySelector("#count");
const body = document.querySelector("#body");
const topic = document.querySelector("#topic");
const composerStatus = document.querySelector("#composer-status");
let sessionReady;

async function loadSession() {
  const response = await fetch("/api/v1/session");
  token = (await response.json()).token;
}

async function loadDiscussions() {
  const response = await fetch("/api/v1/discussions?limit=20");
  if (!response.ok) throw new Error("Could not load discussions");
  const discussions = await response.json() || [];
  count.textContent = discussions.length;
  list.innerHTML = discussions.length ? discussions.map((item) => `<article class="discussion"><div class="discussion-meta"><span class="topic">${escapeHTML(item.topic)}</span><span>${timeAgo(item.createdAt)}</span></div><h3>${escapeHTML(item.body)}</h3><p>Started by ${escapeHTML(item.authorId)}</p></article>`).join("") : '<p class="loading">No conversations yet. Start the first one.</p>';
}

document.querySelector("#composer").addEventListener("submit", async (event) => {
  event.preventDefault();
  const button = event.target.querySelector("button");
  if (!topic.value.trim() || !body.value.trim()) {
    composerStatus.textContent = "Add a topic and your thought first.";
    composerStatus.style.color = "#bd4d32";
    return;
  }
  composerStatus.textContent = "Sending your discussion...";
  composerStatus.style.color = "var(--muted)";
  button.disabled = true;
  try {
    await sessionReady;
    const response = await fetch("/api/v1/discussions", {method:"POST", headers:{"Content-Type":"application/json", Authorization:`Bearer ${token}`}, body:JSON.stringify({topic:topic.value, body:body.value})});
    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || "Could not start the discussion");
    }
    event.target.reset();
    document.querySelector("#chars").textContent = "0 / 500";
    composerStatus.textContent = "Discussion started.";
    composerStatus.style.color = "#4c8060";
    await loadDiscussions();
  } catch (error) {
    composerStatus.textContent = error.message;
    composerStatus.style.color = "#bd4d32";
  } finally {
    button.disabled = false;
  }
});
body.addEventListener("input", () => document.querySelector("#chars").textContent = `${body.value.length} / 500`);
document.querySelector("#refresh").addEventListener("click", loadDiscussions);
function escapeHTML(value) { const div = document.createElement("div"); div.textContent = value; return div.innerHTML; }
function timeAgo(value) { const minutes = Math.max(1, Math.floor((Date.now() - new Date(value)) / 60000)); return minutes < 60 ? `${minutes}m ago` : `${Math.floor(minutes / 60)}h ago`; }
sessionReady = loadSession();
Promise.all([sessionReady, loadDiscussions()]).catch(() => { composerStatus.textContent = "Could not connect to the room."; composerStatus.style.color = "#bd4d32"; list.innerHTML = '<p class="loading">The room is taking a moment. Please refresh.</p>'; });