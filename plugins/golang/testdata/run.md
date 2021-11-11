# Testing Go Run!

some text

<go#run src="cmd" files="main.go,print.go"></go#run>

more text

<go#run src="./cmd" args="-tags,sad" out="stderr"></go#run>

- src - required
- args - optional, defaults to empty
- files - optional, defaults to all files in src
- out - optional, defaults to stdout
