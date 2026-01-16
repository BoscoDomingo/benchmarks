import { arr } from "./setup.ts";

const sum = arr
	.filter((n) => n > 100)
	.filter((n) => n % 2 === 0)
	.map((n) => n * n)
	.reduce((acc, n) => acc + n, 0);

console.log(sum);
