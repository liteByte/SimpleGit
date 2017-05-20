# SimpleGit
A wrapper of git2go, tries to simplify the library by implementing
functions that mimic the most commonly used git commands.

## Install
```
go get github.com/liteByte/SimpleGit
```

You need to have [git2go](https://github.com/libgit2/git2go) in your
GOPATH for SimpleGit to work.

## How to use
The names of the functions are the same as the git commands they mimic,
so what they do is pretty self explanatory if you know git.

```
import "github.com/liteByte/SimpleGit"

// [...]

repo, err := sig.GitClone("https://github.com/liteByte/SimpleGit.git", "master", "projects", sig.MakeGitCredentials("username", "password"))
	if err != nil {
		panic(err)
	}

	// Do stuff

	err = repo.GitAdd("changedFile.txt", "changedFile2.txt")

	author := sig.MakeGitSignature("Name", "Email")
	committer := sig.MakeGitSignature("Name", "Email")

	err = repo.GitCommit("master", "Modify files", author, committer)

	// Don't reuse git credentials, make a new one every time
	// (don't feel like explaining, I'll fix this later)
	err = repo.GitPush("origin", "master", sig.MakeGitCredentials("username", "password"))

```

## Disclaimer
This library is by no means complete, I'll keep adding git commands as I
need them. If you need a command that is missing, create an issue or,
even better, a Pull Request with the command implemented.