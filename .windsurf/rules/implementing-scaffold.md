---
trigger: model_decision
description: Apply when scaffold is mentioned
---

1. Follow user feedback closely, including naming conventions, file organization, and template alignment.
2. When scaffolding new routes or handlers, ensure naming consistency between handler methods, files, and templates (e.g., use underscores, no Handler suffix).
3. Move new handlers to their own file, unless they are different HTTP methods for the same route
4. Align template names with handler files for clarity (e.g., group_by.templ for GroupBy).
5. Register all new routes in the router configuration and update references after renaming handlers.
6. Remove legacy or duplicate files/code after renaming or moving.
7. Add placeholder logic in new handlers for future feature expansion, but keep the code functional.
9. Proactively address user feedback for code organization, clarity, and maintainability throughout the implementation process.
