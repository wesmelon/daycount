import { markComplete } from "./shortcuts.mjs";
import {
  today,
  shift,
  stats,
  addRange,
  validate,
  milestones,
} from "./model.mjs";
const $ = (s) => document.querySelector(s),
  key = "daycount.streaks.v1";
let data = { version: 1, habits: [] },
  selected,
  month = today().slice(0, 7),
  toastTimer,
  storageBlocked = false;
try {
  const raw = localStorage.getItem(key);
  if (raw) data = validate(JSON.parse(raw));
} catch {
  storageBlocked = true;
  setTimeout(
    () =>
      notify(
        "Could not load saved data. Export it from settings before replacing it.",
      ),
    100,
  );
}
selected = data.habits.find(h => !h.archived)?.id;
function notify(s) {
  $("#toast").textContent = s;
  $("#toast").style.display = "block";
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => ($("#toast").style.display = "none"), 5000);
}
function save(next) {
  if (storageBlocked)
    throw Error(
      "Storage could not be loaded. Use settings to export the saved data and restore a valid backup.",
    );
  localStorage.setItem(key, JSON.stringify(next));
  data = next;
  render();
}
function change(fn) {
  const next = structuredClone(data);
  fn(next);
  save(next);
}
function el(tag, cls, text) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (text !== undefined) e.textContent = text;
  return e;
}
function button(text, fn, cls = "") {
  const b = el("button", cls, text);
  b.type = "button";
  b.onclick = () => {
    try {
      fn();
    } catch (e) {
      notify(e.message);
    }
  };
  return b;
}
function habit() {
  return data.habits.find((h) => h.id === selected);
}
function render() {
  const date = today();
  $("#today").textContent = new Date()
    .toLocaleDateString(undefined, {
      weekday: "long",
      month: "long",
      day: "numeric",
    })
    .toUpperCase();
  const active = data.habits.filter(h => !h.archived);
  const done = active.filter((h) => h.entries[date]).length;
  $("#progress").textContent = `${done} of ${active.length} complete`;
  $("#dailyProgress").max = active.length || 1;
  $("#dailyProgress").value = done;
  $("#count").textContent = active.length;
  const list = $("#habits");
  list.replaceChildren();
  if (!active.length) {
    const empty = el("div", "empty");
    empty.append(
      el("h3", "", "Every streak starts with day one."),
      el(
        "p",
        "hint",
        "Create a category like Multivitamin, Reading, or Movement.",
      ),
      button("+ Create your first category", () => category(), "primary"),
    );
    list.append(empty);
  }
  for (const h of active) {
    const s = stats(h.entries),
      card = el("article", `habit ${h.id === selected ? "selected" : ""}`),
      head = el("div", "habit-head");
    head.append(
      el("span", "icon", h.icon),
      button(
        "View calendar ↗",
        () => {
          selected = h.id;
          render();
          $("#detail").scrollIntoView({ behavior: "smooth", block: "start" });
        },
        "details",
      ),
    );
    const menu = el("details", "category-menu");
    const trigger = el("summary", "category-menu-trigger", "⋯");
    trigger.setAttribute("aria-label", "Options for " + h.name);
    const actions = el("div", "category-menu-actions");
    actions.append(
      button("Set up Siri", () => {
        menu.open = false;
        siriSetup(h);
      }),
      button("Archive category", () => {
        menu.open = false;
        change(d => { d.habits.find(x => x.id === h.id).archived = true; });
        selected = data.habits.find(x => !x.archived)?.id;
        render();
        notify("Category archived. Restore it in Backup & settings.");
      }),
    );
    menu.append(trigger, actions);
    menu.addEventListener("toggle", () => {
      if (menu.open) document.querySelectorAll(".category-menu[open]").forEach(other => {
        if (other !== menu) other.open = false;
      });
    });
    head.append(menu);
    card.append(head, el("h3", "", h.name));
    const metric = el("div");
    metric.append(
      el("span", "number", s.current),
      document.createTextNode(s.current === 1 ? " day streak" : " day streak"),
    );
    card.append(
      metric,
      el(
        "p",
        "metric",
        `Personal best ${s.best} days · ${s.total} days recorded`,
      ),
    );
    card.append(
      button(
        h.entries[date] ? "✓ Done today · Undo" : "+ Check in today",
        () => {
          const was = !!h.entries[date],
            best = s.best;
          change((d) => {
            const entry = d.habits.find((x) => x.id === h.id);
            if (was) delete entry.entries[date];
            else entry.entries[date] = { note: "" };
          });
          const now = stats(data.habits.find((x) => x.id === h.id).entries);
          const unlocked = milestones.filter((n) => best < n && now.best >= n);
          notify(
            was
              ? "Today’s check-in removed."
              : unlocked.length
                ? `✦ ${unlocked.at(-1)}-day milestone unlocked!`
                : "One more day in your story. Well done!",
          );
        },
        `check ${h.entries[date] ? "done" : "primary"}`,
      ),
    );

    list.append(card);
  }
  const h = habit();
  $("#detail").hidden = !h || h.archived;
  if (!h || h.archived) return;
  $("#detailName").textContent = h.name;
  const s = stats(h.entries);
  $("#totalDays").textContent = `${s.total} total days`;
  renderCalendar(h);
  $("#rewards").replaceChildren(
    ...milestones.map((n) => {
      const e = el("div", `badge ${s.best >= n ? "unlocked" : ""}`);
      e.append(
        el("b", "", s.best >= n ? "✦" : "◇"),
        el("span", "", `${n} days`),
        el(
          "div",
          "",
          s.best >= n ? "Unlocked" : `${Math.max(0, n - s.current)} to go`,
        ),
      );
      return e;
    }),
  );
  const history = $("#history");
  history.replaceChildren();
  const days = Object.keys(h.entries).sort().reverse();
  if (!days.length)
    history.append(
      el("p", "hint", "Your check-ins and notes will appear here."),
    );
  for (const date of days.slice(0, 60)) {
    const row = el("div", "record"),
      body = el("div");
    body.append(
      el("strong", "", format(date)),
      el("p", "hint", h.entries[date].note || "Checked in"),
    );
    row.append(
      body,
      button("Edit", () => record(date)),
    );
    history.append(row);
  }
  if (days.length > 60)
    history.append(
      el(
        "p",
        "hint",
        "Showing the latest 60 days. Use the calendar to view or edit earlier dates.",
      ),
    );
}
function format(date) {
  return new Date(`${date}T12:00:00`).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}
function renderCalendar(h) {
  const first = `${month}-01`,
    d = new Date(`${first}T12:00:00`),
    end = new Date(d.getFullYear(), d.getMonth() + 1, 0).getDate();
  $("#monthTitle").textContent = d.toLocaleDateString(undefined, {
    month: "long",
    year: "numeric",
  });
  $("#nextMonth").disabled = month >= today().slice(0, 7);
  $("#prevMonth").disabled = month <= "1900-01";
  const grid = $("#calendar");
  grid.replaceChildren();
  for (let i = 0; i < d.getDay(); i++) grid.append(el("span"));
  for (let i = 1; i <= end; i++) {
    const date = `${month}-${String(i).padStart(2, "0")}`,
      b = button(
        String(i),
        () => record(date),
        `${h.entries[date] ? "complete" : ""} ${date === today() ? "today" : ""}`,
      );
    b.disabled = date > today();
    b.setAttribute(
      "aria-label",
      `${format(date)}${h.entries[date] ? ", checked in" : ""}`,
    );
    b.setAttribute("aria-pressed", String(!!h.entries[date]));
    grid.append(b);
  }
}
function moveMonth(n) {
  const d = new Date(`${month}-01T12:00:00`);
  d.setMonth(d.getMonth() + n);
  month = today(d).slice(0, 7);
  render();
}
$("#prevMonth").onclick = () => moveMonth(-1);
$("#nextMonth").onclick = () => moveMonth(1);
function dialog(title) {
  $("#dialogTitle").textContent = title;
  $("#fields").replaceChildren();
  $("#formActions").replaceChildren();
  $("#formError").textContent = "";
  $("#editorForm").onsubmit = (e) => e.preventDefault();
  $("#editor").showModal();
}
function field(label, type, value = "") {
  const l = el("label", "", label),
    i = el(type === "textarea" ? "textarea" : "input");
  if (type !== "textarea") i.type = type;
  i.value = value;
  l.append(i);
  $("#fields").append(l);
  return i;
}
function submit(text, fn) {
  const b = el("button", "primary", text);
  b.type = "submit";
  $("#formActions").append(b);
  $("#editorForm").onsubmit = (e) => {
    e.preventDefault();
    try {
      fn();
      $("#editor").close();
    } catch (err) {
      $("#formError").textContent = err.message;
    }
  };
}
function category(edit = false) {
  const h = edit ? habit() : null;
  dialog(h ? "Edit category" : "A new daily ritual");
  const name = field("Category name", "text", h?.name || "");
  name.required = true;
  name.maxLength = 80;
  name.placeholder = "e.g. Multivitamin";
  const label = el("label", "", "Symbol"),
    select = el("select");
  for (const s of ["✦", "☀", "◉", "♥", "⚡", "☘"]) {
    const o = el("option", "", s);
    o.value = s;
    select.append(o);
  }
  select.value = h?.icon || "✦";
  label.append(select);
  $("#fields").append(label);
  submit(h ? "Save category" : "Create category", () => {
    if (!name.value.trim()) throw Error("Enter a category name.");
    const id = h?.id || crypto.randomUUID();
    change((d) => {
      if (h) {
        const c = d.habits.find((x) => x.id === id);
        c.name = name.value.trim();
        c.icon = select.value;
      } else
        d.habits.push({
          id,
          name: name.value.trim(),
          icon: select.value,
          entries: {},
        });
    });
    selected = id;
    render();
  });
  if (h)
    $("#formActions").append(
      button(
        "Delete category",
        () => {
          if (
            confirm(
              `Delete “${h.name}” and all its records? Export a backup first if you want to keep them.`,
            )
          ) {
            change((d) => (d.habits = d.habits.filter((x) => x.id !== h.id)));
            selected = data.habits[0]?.id;
            render();
            $("#editor").close();
          }
        },
        "danger",
      ),
    );
}
function record(date, range = false) {
  const h = habit();
  dialog(range ? "Add a past streak" : format(date));
  const start = field(range ? "Start date" : "Date", "date", date);
  start.required = true;
  start.min = "1900-01-01";
  start.max = today();
  if (!range) start.disabled = true;
  const end = range ? field("End date (inclusive)", "date", date) : start;
  if (range) {
    end.required = true;
    end.min = "1900-01-01";
    end.max = today();
  }
  const note = field(
    "Note (optional)",
    "textarea",
    range ? "" : h.entries[date]?.note || "",
  );
  note.maxLength = 1000;
  if (range)
    $("#fields").append(
      el(
        "p",
        "hint",
        "Every day in this range will be marked complete. Existing notes in the range will be replaced by this note.",
      ),
    );
  submit(range ? "Save past streak" : "Save check-in", () => {
    change((d) =>
      addRange(
        d.habits.find((x) => x.id === h.id),
        start.value,
        end.value,
        note.value,
      ),
    );
    notify("Your record is saved.");
  });
  if (!range && h.entries[date])
    $("#formActions").append(
      button(
        "Remove check-in",
        () => {
          change(
            (d) => delete d.habits.find((x) => x.id === h.id).entries[date],
          );
          $("#editor").close();
        },
        "danger",
      ),
    );
}
function download(text) {
  const url = URL.createObjectURL(
    new Blob([text], { type: "application/json" }),
  );
  const a = el("a");
  a.href = url;
  a.download = `daycount-${today()}.json`;
  a.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
$("#backup").onclick = () => {
  dialog("Your data, in your hands");
  $("#fields").append(
    el(
      "p",
      "hint",
      "Records live in this browser on this device. They do not sync between devices. Clearing website data removes them. Export regularly and before switching browsers or adding the app to your Home Screen.",
    ),
    el(
      "p",
      "hint",
      "On iPhone: open the hosted app in Safari, tap Share, then Add to Home Screen. After the first visit over HTTPS, the app also works offline.",
    ),
  );
  $("#formActions").append(
    button("Export backup", () =>
      download(localStorage.getItem(key) || JSON.stringify(data)),
    ),
  );
  const archived = data.habits.filter(h => h.archived);
  $("#fields").append(el("h3","","Archived categories"));
  if (!archived.length) $("#fields").append(el("p","hint","No archived categories."));
  for (const h of archived) {
    $("#fields").append(button("Restore " + h.name, () => {
      change(d => { d.habits.find(x=>x.id===h.id).archived=false; });
      selected=h.id; render(); $("#editor").close(); notify("Category restored.");
    },"wide"));
  }
  const input = field(
    "Restore a JSON backup (replaces all categories)",
    "file",
  );
  input.accept = ".json,application/json";
  input.onchange = async () => {
    try {
      const file = input.files[0];
      if (!file) return;
      if (file.size > 10000000) throw Error("Backup must be under 10 MB.");
      const next = validate(JSON.parse(await file.text()));
      if (
        !confirm(
          `Replace all current records with ${next.habits.length} categories from this backup?`,
        )
      )
        return;
      localStorage.setItem(key, JSON.stringify(next));
      storageBlocked = false;
      data = next;
      selected = data.habits[0]?.id;
      render();
      $("#editor").close();
      notify("Backup restored.");
    } catch (e) {
      $("#formError").textContent = e.message;
    }
  };
};
$("#closeDialog").onclick = () => $("#editor").close();
$("#newHabit").onclick = () => category();
$("#editHabit").onclick = () => category(true);
$("#addHistory").onclick = () => record(today(), true);
window.addEventListener("storage", (e) => {
  if (e.key === key) {
    try {
      data = e.newValue
        ? validate(JSON.parse(e.newValue))
        : { version: 1, habits: [] };
      storageBlocked = false;
      if (!habit()) selected = data.habits[0]?.id;
      render();
    } catch {
      storageBlocked = true;
      notify("Saved data could not be read.");
    }
  }
});
let lastDate = today();
setInterval(() => {
  if (today() !== lastDate) {
    lastDate = today();
    render();
  }
}, 30000);
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) render();
});
render();
if ("serviceWorker" in navigator)
  navigator.serviceWorker.register("./sw.js").catch(() => {});

function siriSetup(h) {
  dialog("Set up Siri");
  const url = new URL(location.pathname, location.origin);
  url.hash = new URLSearchParams({checkin:h.id}).toString();
  $("#fields").append(el("p","hint",'In Apple Shortcuts, create a shortcut named “Log '+h.name+'”. Add a URL action with the link below, then an Open URLs action. Say “Hey Siri, Log '+h.name+'”.'));
  const input = field("Check-in link","text",url.href);
  input.readOnly = true;
  input.onclick = () => input.select();
  $("#fields").append(el("p","hint","Use this with the browser that holds your records. Safari and a Home Screen app may have separate storage. If the category is missing, export your records from the original app and restore them in the browser opened by the shortcut. You may need to unlock your iPhone."));
  $("#formActions").append(button("Copy link", () => {
    navigator.clipboard.writeText(url.href).then(()=>notify("Link copied.")).catch(()=>{input.select();$("#formError").textContent="Select and copy the link above.";});
  },"primary"));
}
function handleCheckinLink() {
  const params = new URLSearchParams(location.hash.slice(1));
  if (!params.has("checkin")) return;
  const id = params.get("checkin");
  history.replaceState(null,"",location.pathname+location.search);
  try {
    if (storageBlocked) throw Error("Saved records could not be loaded. Restore a backup before checking in.");
    // Read fresh storage when an existing browser tab receives another shortcut.
    const raw = localStorage.getItem(key);
    const next = raw ? validate(JSON.parse(raw)) : {version:1,habits:[]};
    const result = markComplete(next,id,today());
    selected = id;
    save(next);
    notify(result.name + (result.already ? " is already complete today." : " checked in for today!"));
  } catch(error) {
    dialog("Check-in could not be saved");
    $("#fields").append(el("p","hint",error.message),el("p","hint","Open this link in the browser containing your records. If you use the Home Screen app, export a backup there and restore it here, then run the shortcut again."));
  }
}
window.addEventListener("hashchange",handleCheckinLink);
handleCheckinLink();

document.addEventListener("click", event => {
  document.querySelectorAll(".category-menu[open]").forEach(menu => {
    if (!menu.contains(event.target)) menu.open = false;
  });
});
document.addEventListener("keydown", event => {
  if (event.key !== "Escape") return;
  document.querySelectorAll(".category-menu[open]").forEach(menu => {
    menu.open = false;
    menu.querySelector("summary").focus();
  });
});
