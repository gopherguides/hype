# Testing Go Run!

- exec: command to run - required
- dir: directory to run the command in, defaults to current directory
- hide-stdout: hide stdout output, defaults to false
- hide-stderr: hide stderr output, defaults to false
- hide-cmd: hide the command used to run the command, defaults to false

more text

<cmd exec="go run -tags sad ." src="./cmd">

As we can see from the output, the command failed.

</cmd>

even more output!

<cmd exec="tree" src="./cmd"></cmd>

clean the pwd

<cmd exec="echo hello" src="./cmd"></cmd>

skip stderr

<cmd exec="go run main.go print.go" src="./cmd" hide-stdout></cmd>

<cmd exec="go run main.go print.go" src="./cmd" hide-cmd>

As we can see from the output, the command was executed successfully.

</cmd>

<include src="sub/sub.md"></include>
