# Week 1 - Getting Started with Go - TODO

overview

<include src="basics.md"></include>

<include src="strings.md"></include>

<include src="utf8.md"></include>

<include src="numbers.md"></include>

<include src="booleans.md"></include>

<include src="variables.md"></include>

<include src="constants.md"></include>

# Exercise (Due Wednesday)

Locally create a new folder named, "gopherguides-intro-to-go". Inside of the folder initialize a new git repository. Next, create a new folder named `gopherguides-intro-to-go/week01`. Inside of the `week01` folder initialize a new Go module named "github.com/YOUR-USERNAME/gopherguides-intro-to-go/week01".

```
$ go mod init github.com/YOUR-USERNAME/gopherguides-intro-to-go/week01
```

Commit the `go.mod` to the git repository.  Next, on [GitHub.com](https://github.com/) create a new public repository named "gopherguides-intro-to-go" under your account and upload your local repository following the instructions on GitHub.

```text
.
└── gopherguides-intro-to-go
    └── week01
        └── go.mod
```

---

# Assignment 1 (Due Sunday)

## 1.1

Write a "Hello, World" style Go program using the `main` package. Your file should be named `main.go`. This program **must** compile and print "Hello, World!", with a new line after it, to the console window when run. Publish this code to your repository under your `week01` folder you created earlier this week. Next, create a branch in your local project called, `assignment01`. Using [pkg.go.dev](pkg.go.dev) research the `fmt` package. Use the `fmt` package to print "Printing, TODO!", replacing "TODO" with the proper printing verb to properly print the following types: `string ("Go"), int (42), bool (true)`. Use `go vet` to confirm you are using the correct verb. Finally, open a PR to merge your new changes into your `main` branch. This PR should contain a paragraph or two explaining the changes and how they were implemented. Submit the link to the PR to be reviewed.

## 1.2

Write a short essay describing your history in technology and how you feel that Go fits into your plans for your future. Additionally, write a short essay discussing any surprises you found when researching the `fmt` package. Include how does printing in Go differ from other languages you may have used before. Please be specific and cite examples. (500 words minimum)
