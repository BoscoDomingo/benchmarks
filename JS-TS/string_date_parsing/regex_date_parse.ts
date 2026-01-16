import { cleanDateString } from "./setup.ts";

const DATE_FORMAT_REGEX = /^(\d{2})\/(\d{2})\/(\d{4})$/; // DD/MM/YYYY

const match = cleanDateString.match(DATE_FORMAT_REGEX);
if (!match) {
	throw new Error(
		`Invalid date format "DD/MM/YYYY" expected. inputDate received: ${cleanDateString}`,
	);
}

const dayOfMonth = match[1];
const month = match[2];
const year = match[3];

const result = new Date(Date.parse(`${year}-${month}-${dayOfMonth}`));

console.log(result);
