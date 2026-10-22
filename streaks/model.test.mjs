import { test } from "node:test";
import assert from "node:assert/strict";
import {
  stats,
  addRange,
  shift,
  validate,
  validDate,
  today,
} from "./model.mjs";
test("local date and DST-safe calendar arithmetic", () => {
  assert.equal(today(new Date(2026, 0, 2, 23)), "2026-01-02");
  assert.equal(shift("2024-02-28", 1), "2024-02-29");
  assert.equal(shift("2026-03-08", 1), "2026-03-09");
  assert.equal(validDate("2026-02-29"), false);
});
test("current streak allows today to remain incomplete, but not yesterday", () => {
  const h = { entries: {} };
  addRange(h, "2026-01-01", "2026-01-10", "note", "2026-01-20");
  assert.deepEqual(stats(h.entries, "2026-01-11"), {
    current: 10,
    best: 10,
    total: 10,
  });
  assert.equal(stats(h.entries, "2026-01-12").current, 0);
  h.entries["2026-01-12"] = { note: "" };
  assert.deepEqual(stats(h.entries, "2026-01-12"), {
    current: 1,
    best: 10,
    total: 11,
  });
});
test("overlap deduplicates, notes update, deletion splits streak", () => {
  const h = { entries: {} };
  addRange(h, "2026-01-01", "2026-01-10", "a", "2026-02-01");
  addRange(h, "2026-01-05", "2026-01-12", "b", "2026-02-01");
  assert.equal(stats(h.entries, "2026-01-12").total, 12);
  assert.equal(h.entries["2026-01-05"].note, "b");
  delete h.entries["2026-01-06"];
  assert.equal(stats(h.entries, "2026-01-12").best, 6);
});
test("reject invalid or future ranges", () => {
  for (const [a, b] of [
    ["2026-02-30", "2026-03-01"],
    ["2026-01-10", "2026-01-01"],
    ["2026-01-01", "2026-02-02"],
  ])
    assert.throws(() => addRange({ entries: {} }, a, b, "", "2026-02-01"));
});
test("backup round trip and hostile malformed inputs", () => {
  const data = {
    version: 1,
    habits: [
      {
        id: "a",
        name: "Vitamin",
        icon: "✦",
        entries: { "2026-01-01": { note: "Hello" } },
      },
    ],
  };
  assert.deepEqual(validate(JSON.parse(JSON.stringify(data))), data);
  assert.throws(() =>
    validate({ ...data, habits: [...data.habits, ...data.habits] }),
  );
  assert.throws(() =>
    validate({
      version: 1,
      habits: [{ ...data.habits[0], entries: { oops: { note: "" } } }],
    }),
  );
});
