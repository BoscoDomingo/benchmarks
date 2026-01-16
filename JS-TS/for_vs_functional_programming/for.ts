import { arr } from "./setup.ts";

let sum = 0;

for (const n of arr) {
	if (n <= 100) {
		continue;
	}

	if (n % 2 !== 0) {
		continue;
	}

	sum += n * n;
}

console.log(sum);
