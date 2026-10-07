###History file

	- This file is intended to include the detailed task planning, WITH TRACKING STATUSES.
	- Here is the space for keeping track of the project's progress in the loop described in 
	/ai/AGENTS.md

###TASKS QUEUE

[
...
TASK 001 [DONE] - Create a new button on 'x' > 'y' form that does yada yada...
		- Summ: - Button created, we invoke 'dummyFunction' that triggers a 'dummy-event'...
TASK 002 [REVIEW] - Create a 'Refresh Token' button on 'b' > 'c' list...
	- ERROR: We tried to use dummyLib to produce the dummyFeature but this lib lacks support
	for dummyUtil, operator rejected this implementation.
TASK 003 [PENDING]- Do awesome button on 'a' > 'b'
...
] 

### Local quality gate — 2026-10-07

- Added local `make lint`, `make complexity`, `make duplication`, and canonical
  `make check`, reusing the full existing tests and native builds. No CI added.
- Pinned golangci-lint setup and frontend ESLint/TypeScript/Vue/jscpd dependencies;
  excluded Ent, frontend dependencies, generated code, coverage, and build outputs.
- Recorded pre-existing Go complexity/cleanup-error debt with narrow exceptions
  and existing frontend template duplication in `VALIDATION.md`.
- Updated `AGENT.md` and `start.md`: every implementation must end with a passing
  root `make check`; fix failures and repeat until clean.
- Individual checks, full tests/builds, exclusion fixtures, deliberate failure
  propagation probes, and fresh `npm ci` passed. See `VALIDATION.md` for the final
  complete-gate result. The placeholder task queue and next-task template are retained.
