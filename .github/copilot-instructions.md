# copilot-instructions.md

## General Guidelines
- Follow Go best practices and idiomatic code style (Effective Go).
- Use clear, concise, and self-documenting code.
- Avoid unnecessary complexity; prefer readability over clever tricks.
- Always handle errors explicitly unless clearly safe to ignore.
- Prefer dependency injection over global variables.
- **All chat responses must be in Japanese**, even if code comments or examples contain English.
- **After making any changes or implementing new code, always run the tests and ensure all of them pass before finalizing.**

---

## goDoc Updates
- Ensure all exported functions, methods, types, constants, and variables have a GoDoc comment.
- GoDoc comments must start with the name of the item they describe.
- Explain **what** the function/type does, not how it is implemented (unless relevant for correct usage).
- Example:
  ```go
  // CalculateTax returns the tax amount for the given income.
  ```

---

# Testing
- Do not run `go test` directly. Watch 
- When run tests, always use the VS Code task labeled `unit test`.