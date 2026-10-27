import {test} from 'node:test';
import assert from 'node:assert/strict';
import {markComplete} from './shortcuts.mjs';
test('check-in is idempotent and preserves notes',()=>{
const data={habits:[{id:'a',name:'Vitamin',entries:{}}]};
assert.equal(markComplete(data,'a','2026-10-26').already,false);
data.habits[0].entries['2026-10-26'].note='with lunch';
assert.equal(markComplete(data,'a','2026-10-26').already,true);
assert.equal(data.habits[0].entries['2026-10-26'].note,'with lunch');
assert.equal(Object.keys(data.habits[0].entries).length,1);
markComplete(data,'a','2026-10-27');
assert.equal(Object.keys(data.habits[0].entries).length,2);
});
test('unknown category cannot create or modify records',()=>{
const data={habits:[{id:'a',name:'Vitamin',entries:{}}]};
assert.throws(()=>markComplete(data,'missing','2026-10-26'));
assert.deepEqual(data.habits[0].entries,{});
});
test('archived category cannot receive Siri check-in',()=>{const data={habits:[{id:'a',name:'Vitamin',archived:true,entries:{}}]};assert.throws(()=>markComplete(data,'a','2026-10-26'));assert.deepEqual(data.habits[0].entries,{});});