###Planning file

	- This file is intended to include the detailed file planning, divided by tasks,
	for the ongoing project.
	- This file currently is empty in initial state, and will be filled once in the initial 
	prompt with the following sections:

###SUMMARY

[ ... include a detailed summary with tech stack, folders structure, archtechture,
  featured libraries and project objective here.
]

###STANDARDS

[
... include details about defined standards trought the code. This section can be tweaked when
explictly address a approach change on a specific standard
]

Examples:

	- Project uses MVC architechture, so controllers are responsible for handling async
trycatch lifecycle. Services embrace core logic like performing internal DB operations, and external
API calls (...)
	- All forms in the application follow a standardized request lifecycle using 
DummyFormLib and DummyReqLib, leveraging dummy.method() mentioned att dummy.docs.io (...)


###PLANNING

 [ This section includes the detailed context for each task before it starts, without task status,
 and define pratical guidelines for how we should implement each task]

Examples:

[
...
TASK 001 - List x using xAPI that exposes endpoint GET '/x/api/v1/dummy-endpoint.'.
	The body is 
req.dummy
{
yada
} (...)
The official documentation states that token 'abc' should be used for a process. In our backend,
we will perform operation yada on yadayada service (...)
	
TASK 002 - Do y
TASK 003 - Do awesome button on 'a' > 'b'
..
] 
