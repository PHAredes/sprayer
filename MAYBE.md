# MAYBE.md - Project Evolution (Functional & Concise)

This document outlines potential improvements for Sprayer, focusing on functional programming patterns, concise naming, and architectural refinement.

## 1. Internalize `parse` into `job`
*   **Reasoning**: `parse` is almost always used in the context of job data extraction.
*   **Pros**: Flattens the package hierarchy.
*   **Cons**: Slightly larger `job` package.

## 2. Refactor Scrapers to `Stream[Job]`
*   **Reasoning**: Use channels as streams/observables rather than returning slices.
*   **Pros**: Lower memory footprint, immediate UI feedback.
*   **Cons**: More complex error handling across the stream.

## 3. Replace Structs with Maps for Metadata
*   **Reasoning**: Use a `Map` for optional job metadata (salary, tech stack) instead of fixed struct fields.
*   **Pros**: Flexibility for different job boards.
*   **Cons**: Loss of type safety for specific fields.

## 4. Rename `internal/scraper` to `internal/src`
*   **Reasoning**: "Scraper" is a long word. "Src" is concise and standard.
*   **Pros**: Shorter imports.
*   **Cons**: "Src" might be confused with general source code.

## 5. Rename `internal/profile` to `internal/me`
*   **Reasoning**: "Profile" is corporate. "Me" is concise and personal.
*   **Pros**: Very clear intent.
*   **Cons**: Might feel too informal for some.

## 6. Functional Scorer Pipeline
*   **Reasoning**: Define scoring as `type Scorer func(Job) int`. Chain them with `Sum(scorers...)`.
*   **Pros**: Easy to add/remove scoring criteria (e.g., tech match, location match).
*   **Cons**: Slightly more overhead than a single function.

## 7. Immutable Job States
*   **Reasoning**: Instead of mutating a `Job` struct, return a new copy on every pipeline step.
*   **Pros**: Thread safety, easier debugging.
*   **Cons**: More allocations (mitigated by Go's efficient stack/heap management).

## 8. Rename `Apply` to `Do`
*   **Reasoning**: "Apply" is specific to jobs. "Do" is the action of the agent.
*   **Pros**: Minimalist.
*   **Cons**: Potentially too vague.

## 9. Monadic Error Handling (Result Type)
*   **Reasoning**: Implement a simple `Result[T]` struct to avoid `if err != nil` everywhere.
*   **Pros**: Cleaner functional chains.
*   **Cons**: Not idiomatic Go; can be clunky without proper generics support in all libs.

## 10. Curried Scrapers
*   **Reasoning**: Scrapers should be functions that take config and return the scraping function.
*   **Pros**: Better testability through dependency injection.
*   **Cons**: Deeply nested function signatures.

## 11. Replace `internal/apply` with `internal/out`
*   **Reasoning**: Everything going "out" of the system (email, export, logs).
*   **Pros**: Consolidates output logic.
*   **Cons**: "Out" is a reserved word in some contexts (not Go, but common).

## 12. Use `Context` for Pipeline Cancellation
*   **Reasoning**: Pass `context.Context` through all functional chains.
*   **Pros**: Proper timeout handling for slow scrapers or LLM calls.
*   **Cons**: More boilerplate in function signatures.

## 13. Declarative TUI (Functional UI)
*   **Reasoning**: Move away from stateful Bubble Tea models towards a more "view as a function of state" approach.
*   **Pros**: Easier to reason about UI transitions.
*   **Cons**: Bubble Tea is inherently designed around the Model-Update-View loop.

## 14. Global `Store` as an Interface
*   **Reasoning**: Define `type DB interface`.
*   **Pros**: Swap SQLite for Postgres or In-memory easily.
*   **Cons**: Adds a layer of abstraction.

## 15. Short-lived Scraper Processes
*   **Reasoning**: Run browser scrapers (Rod) in separate processes to prevent memory leaks.
*   **Pros**: Robustness.
*   **Cons**: Higher startup latency for each scrape.

## 16. Use `λ` (Lambda) for Internal Variable Names
*   **Reasoning**: Use short names like `j` for job, `p` for profile, `s` for store in local scopes.
*   **Pros**: Reduces visual noise.
*   **Cons**: Can be cryptic for new contributors.

## 17. Pattern-Based Deduplication
*   **Reasoning**: Instead of ID-only dedup, use fuzzy hashing on job descriptions.
*   **Pros**: Catches the same job posted on different boards.
*   **Cons**: Computationally expensive.

## 18. Consolidate `cmd/` into a single `main.go`
*   **Reasoning**: Use a single binary with subcommands (like `git`).
*   **Pros**: Easier distribution.
*   **Cons**: Larger binary size if some modes aren't used.

## 19. Functional Prompt Templates
*   **Reasoning**: Prompts as Go functions instead of `.txt` files.
*   **Pros**: Compile-time safety for variables.
*   **Cons**: Harder for non-coders to edit prompts.

## 20. Rename `internal/llm` to `internal/ai`
*   **Reasoning**: "AI" is the common term; "LLM" is technical jargon.
*   **Pros**: Shorter and more recognizable.
*   **Cons**: "AI" is a very broad term.
