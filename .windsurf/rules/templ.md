---
trigger: model_decision
description: Check this rule when writing in ".templ" files
---

# templ Syntax & Best Practices Rules

This document summarizes key rules and discoveries for using the `templ` language in Go web projects:

## 1. Go Code Blocks
- All Go code (variable declarations, logic, loops, etc.) must be inside `{{ ... }}` blocks.
- Do not mix Go code and HTML nodes inside the same `{{ ... }}` block.
- Close the Go code block (`}}`) before writing any HTML or template markup.

## 2. HTML and Template Markup
- All HTML nodes and template markup must be outside of `{{ ... }}` Go code blocks.
- Use template loops and expressions like `for ... { }` directly in the HTML section.
- Do not declare variables or run logic in the HTML section—do this in a Go block above.

## 3. Variable Scope
- Variables declared in a `{{ ... }}` block are accessible in subsequent template code, including HTML and template for-loops.
- Avoid redeclaring variables with `:=` if they are already declared in the outer scope—use `=` to assign instead.

## 4. Layouts and Nodes
- `@ui.Layout(nil)` (or any layout) must wrap only HTML nodes and template markup, not Go code blocks.
- Do not wrap a `{{ ... }}` block around a layout directive or other keywords like if `for` and `if`.
- Do not wrap a `{{ ... }}` block around the for loop in {{ ... }}. The templ engine parses Go statements like for natively in the template body.
- The layout must contain at least one HTML node or valid template expression.

## 5. Duplication and Structure
- Only perform grouping, mapping, or other logic once at the top of the `templ` declaration.
- Do not duplicate Go code blocks or logic after the HTML section.
- Keep the template structure: logic at the top, then layout and HTML rendering.
- If there is a handler that is using this template, consider moving the logic to the handler and passing the results to the template.

## 6. Error Avoidance
- Common errors include:
  - `expected nodes, but none were found` (usually caused by Go code blocks inside layouts or missing HTML nodes)
  - `missing close braces` (unclosed `{{ ... }}` blocks)
  - `no new variables on left side of :=` (variable redeclaration)

## Example Structure
```templ
{{
    // Go logic here
    groupMap := map[string][]Email{}
    groupKeys := []string{}
    // ...
}}

@ui.Layout(nil) {
    <h1>Title</h1>
    for _, group := range groupKeys {
        <!-- HTML rendering here -->
    }
}
```
3. Template Loop Syntax
Old (Incorrect):
```bad-templ
{{ for _, group := range groupKeys { }}
    ...
    {{ for _, email := range groupMap[group] { }}
        ...
    {{ } }}
{{ } }}
```

New (Correct):
```templ
for _, group := range groupKeys {
    ...
    for _, email := range groupMap[group] {
        ...
    }
}

```
You switched from the templ curly-brace loop syntax ({{ ... }}) to direct Go-style loop syntax (for _, ... { ... }). This matches the templ rules for embedding Go code directly in the template body, which is required for proper code generation.
---
These rules are based on discoveries and debugging from implementing advanced grouping logic in Gmail Sorter templates. Follow them to avoid common templ compilation and runtime errors.
