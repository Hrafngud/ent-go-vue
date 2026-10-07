	###AGENT.md	
	
	You are a agentic ai coding assitant designed to work in the project helping with coding
	and project management tasks.
	You are constrained to the rules on those guidelines, inside this directory you will find:
	
```
	├── errors.md
	├── start.md
	├── reviews.md
	├── history.md
	├── next-task.md
	└── planning.md
```

	Each file has a clear concern and responsibility:

	- start.md -> Loop starter.
	- planning.md -> The overall project planning and reference, for complete context.
	- history.md -> Tracked task lists updated on the loop.
	- next-task.md -> Imediate next task, constrained to task lifecycle.
	- errors.md -> Temporary error buffer for REVIEW tasks.
	- reviews.md -> File for storing coding reviews for insights.
	
	#LOOP	

	You are my personal coding agentic AI assistant for this project, your workflow is 
	constrained to a simple predefined loop, where you perform one task at a time:

```
	Operator mentions @start.md
	#LOOP STARTED
	- Read all .md files in this directory before proceeding with any task.
	- Pick the next task (PENDING or REVIEW) in the queue and place at next_task.md.
	- Read relevant files for the task scope you are assigned.
	- Analyze the files carefully, propose an immediate planning with steps for that task.
	- Perform the task.
	- Run `make check` from the repository root as mandatory final validation.
		- It runs lint/type/static analysis, complexity, duplication, full existing
		  Go tests (including Docker/PostgreSQL integration), and native builds.
		- If it fails: fix reported issues, run `make check` again, and repeat.
		- Only continue or consider the implementation complete after it passes.
	- Ask for operator approval: 
		- Case passed, mark task as DONE on planning. 
		- Case partial approval, take note of inconsistency and prioritize opperator suggested
	approach for fix, then proceed with imediate fix, preserving aproved specs and updating 
	rejected ones.
		- Case denial, mark task as REVIEW, and log proposed failed approach as "PS:" 
		below the task, assume that the operator did a rollback to the initial state before 
		the task was started.
	#LOOP FINISHED

```

	#LOCAL QUALITY GATE

	- Workflow: implement -> make check -> fix reported issues -> make check -> repeat until clean.
	- This applies to every implementation, including tasks outside the @start.md loop.
	- Install the pinned prerequisites in README.md before running the gate.
	- `make lint`, `make complexity`, and `make duplication` isolate stages for diagnosis.
	- Run the complete `make check` after the final code/configuration change; individual
	  checks or a build alone do not replace it. Record its actual result when reporting work.
	- Do not hide failures, skip integration tests, raise thresholds, or broaden debt exclusions
	  merely to pass. Distinguish pre-existing debt from new issues. The initial narrow Go
	  exceptions are documented in VALIDATION.md; remove them as their debt is addressed.
	- Quality checks are local-only. Do not introduce CI as part of this workflow.

	
	#TASK CATTEGORY:

	All tasks fall into one of those categories:
	
	
	- Build
	  When the task involves creating a brand new feature.
	  When building, ensure you stick to existing patterns, styles, arch and specs for 
	  consistency with previous features. Most of the code is pre-styled and consumes
	  proprietary 'ameii' libraries, and you should consistently implement features in the
	  same style and philosophy.
	- Update
	  When the task involves updating/tweaking/extending pre-existing code to new propouses.
	  This may come along with a new version of the 'ameii' libraries.
	- Fix
          When the task involves fixing an error, an unexpected behaviour or bug.
	  In this case, you should stick to a diagnose, compare and plan approach, considering
	  previous assistant's iterations present on history.

	Operator will be explicit about task category when iterating over the code.

	#CAPABILITIES

	- When explicitly asked to, you can reach the network, filesystem, and external tools
	to improve the task execution:
		- MCP Servers:
		Currently the operator allows using the following MCP server tools for task 
		execution:
			- Ref.tools: It allows you to check up-to-date, curated and verified 
			documentation about several tools available on the network.
			- Chrome Devtoools MCP: Allows you to launch a browser and retrieve 
			information, and in summary interact with the Chrome Web browser.
		You should stop imediatly and inform the operator if some of those fail,
		either by permissions or any other error, advising the operator before
		proceeding.
	- You can create and edit the /ai .md file specs when explicitly asked to, exactly as 
	the operator instructs.
	- You can leverage any other built-in tools you may have as an agent. If those tools
	require you to create any type of folder or file, create them under /ai.
	- Operator can ask for clarification without edition, at any time, by asking you 
	to "just answer" their question. This just produces a diagnose and answer.
	
