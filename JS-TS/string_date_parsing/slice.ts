import { cleanDateString } from "./setup.ts";

const dayOfMonth = cleanDateString.slice(0, 2);
const month = cleanDateString.slice(3, 5);
const year = cleanDateString.slice(6, 10);

const result = new Date(Date.parse(`${year}-${month}-${dayOfMonth}`));

console.log(result);
