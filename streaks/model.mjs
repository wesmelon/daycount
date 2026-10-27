export const milestones = [10, 30, 50, 100, 365];
export function today(now = new Date()) {
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
}
export function validDate(s) {
  if (typeof s !== "string" || !/^\d{4}-\d{2}-\d{2}$/.test(s)) return false;
  const d = new Date(`${s}T12:00:00Z`);
  return !isNaN(d) && d.toISOString().slice(0, 10) === s && s >= "1900-01-01";
}
export function shift(s, n) {
  const d = new Date(`${s}T12:00:00Z`);
  d.setUTCDate(d.getUTCDate() + n);
  return d.toISOString().slice(0, 10);
}
export function stats(entries, date = today()) {
  const days = Object.keys(entries)
    .filter((d) => validDate(d) && d <= date)
    .sort();
  let best = 0,
    run = 0,
    previous;
  for (const d of days) {
    run = previous && shift(previous, 1) === d ? run + 1 : 1;
    best = Math.max(best, run);
    previous = d;
  }
  let current = 0,
    cursor = entries[date] !== undefined ? date : shift(date, -1);
  while (entries[cursor] !== undefined) {
    current++;
    cursor = shift(cursor, -1);
  }
  return { current, best, total: days.length };
}
export function addRange(habit, start, end, note, date = today()) {
  if (!validDate(start) || !validDate(end) || start > end || end > date)
    throw Error("Choose valid dates, ending today or earlier.");
  if ((new Date(end) - new Date(start)) / 86400000 > 36600)
    throw Error("Please enter a range of 100 years or less.");
  for (let d = start; d <= end; d = shift(d, 1))
    habit.entries[d] = { note: note.trim().slice(0, 1000) };
}
export function validate(data) {
  if (
    !data ||
    data.version !== 1 ||
    !Array.isArray(data.habits) ||
    data.habits.length > 200
  )
    throw Error("This is not a valid Daycount backup.");
  const ids = new Set();
  return {
    version: 1,
    habits: data.habits.map((h) => {
      if (
        !h ||
        typeof h.id !== "string" ||
        !h.id ||
        ids.has(h.id) ||
        typeof h.name !== "string" ||
        !h.name.trim() ||
        h.name.length > 80 ||
        !h.entries ||
        Array.isArray(h.entries) ||
        typeof h.entries !== "object"
      )
        throw Error("Invalid category in backup.");
      ids.add(h.id);
      const entries = {};
      for (const [d, e] of Object.entries(h.entries)) {
        if (
          !validDate(d) ||
          d > today() ||
          !e ||
          typeof e.note !== "string" ||
          e.note.length > 1000
        )
          throw Error("Invalid date or note in backup.");
        entries[d] = { note: e.note };
      }
      return {
        id: h.id,
        name: h.name.trim(),
        icon: ["✦", "☀", "◉", "♥", "⚡", "☘"].includes(h.icon) ? h.icon : "✦",
        ...(h.archived === true ? {archived:true} : {}),
        entries,
      };
    }),
  };
}
