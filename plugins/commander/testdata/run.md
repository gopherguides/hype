# Testing Go Run!

- exec: command to run - required
- dir: directory to run the command in, defaults to current directory
- hide-stdout: hide stdout output, defaults to false
- hide-stderr: hide stderr output, defaults to false
- hide-cmd: hide the command used to run the command, defaults to false

<cmd exec="go run main.go print.go" dir="./cmd" hide-cmd>

As we can see from the output, the command was executed successfully.

</cmd>

more text

<cmd exec="go run -tags sad ." dir="./cmd">

As we can see from the output, the command failed.

</cmd>

even more output!

<cmd exec="tree" dir="./cmd"></cmd>

clean the pwd

<cmd exec="echo hello" dir="./cmd"></cmd>

skip stderr

<cmd exec="go run main.go print.go" dir="./cmd" hide-stdout></cmd>
