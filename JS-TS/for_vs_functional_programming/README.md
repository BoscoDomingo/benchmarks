# For Loop vs Functional Programming vs Iterators

## Use Case

Computes the **sum of squares of even numbers greater than 100** from an array of 5000 random integers (1-1000).

This operation involves:
- **Filtering** (n > 100)
- **Filtering** (n is even)
- **Mapping** (n → n²)
- **Reducing** (sum)

## Approaches

| File          | Method                                       | Passes | Intermediate Arrays |
| ------------- | -------------------------------------------- | ------ | ------------------- |
| `for.ts`      | Index-based for loop                         | 1      | 0                   |
| `fp.ts`       | `.filter().filter().map().reduce()`          | 4      | 3                   |
| `iterator.ts` | `.values().filter().filter().map().reduce()` | 1      | 0                   |

## Key Differences

- **For loop**: Most performant, single pass, manual state management
- **FP methods**: Most readable, but creates intermediate arrays on each chained call
- **Iterator helpers**: FP-like API with lazy evaluation (single pass at consumption). Uses native Iterator Helpers (ES2025+)
