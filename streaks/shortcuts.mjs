export function markComplete(data, id, date) {
  const habit = data.habits.find(h => h.id === id);
  if (!habit) throw Error("This category was not found in this browser. No check-in was added.");
  if (habit.archived) throw Error("This category is archived. Restore it in Backup & settings first.");
  const already = Object.hasOwn(habit.entries,date);
  if (!already) habit.entries[date] = {note:""};
  return {name:habit.name,already};
}
