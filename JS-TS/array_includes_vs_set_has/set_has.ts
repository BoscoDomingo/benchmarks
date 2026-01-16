import { randomArray } from "./setup.ts";

const set = new Set(randomArray);
const randomValueInSet = set.has(Math.random());

console.log(randomValueInSet);
