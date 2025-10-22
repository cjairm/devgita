## 🧭 Why

Ensure consistent, user-friendly, and predictable error handling across all commands using `utils.MaybeExitWithError()`.
Goals:

- Centralize CLI exit behavior
- Provide clear messages and uniform exit codes
- Support verbose/debug modes
- Simplify rollback and recovery logic

---

## ⚙️ Error Handling Patterns

### 1. Return, Don’t Panic

Always return errors to the caller instead of panicking. Panics are reserved for truly unrecoverable conditions.

### 2. Add Context When Propagating

Each layer should wrap or annotate errors to clarify where and why they occurred.

### 3. Centralize CLI Exit Logic

Use a single utility to handle program termination and message formatting, instead of scattering `os.Exit` or print statements.

### 4. Rollback on Failure

Design functions to safely revert partial operations when an error occurs. Use deferred cleanup steps triggered only if an error is returned.

### 5. Integrate Verbose Logging

Allow a `--verbose` or `--debug` flag to print detailed error information while keeping user-facing messages short and friendly by default.

### 6. Use Structured User Errors

Differentiate between internal errors (for debugging) and user-friendly messages (for CLI output). Include error codes where relevant.

### 7. Classify and Inspect Errors

Define well-known error types or sentinel values for predictable control flow (e.g., “not found,” “already exists”).

### 8. Keep Messages Friendly

Avoid exposing raw technical details to end users. Provide actionable, human-readable error messages instead.

---

✅ **Summary**

| Goal         | Pattern                    |
| ------------ | -------------------------- |
| Consistency  | Centralized error handling |
| Clarity      | Context-rich propagation   |
| Reliability  | Safe rollback on failure   |
| Transparency | Optional verbose mode      |
| Usability    | Human-friendly messages    |
| Control      | Standardized exit codes    |
