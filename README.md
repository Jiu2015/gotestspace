## gotestspace

gotestspace is used to quickly create a working directory for shell execution using go, as well as a tool for customizing the execution of the shell. It can help you quickly create an independent workspace for unit testing, improve the efficiency of unit test writing, and improve the happiness of developers.		

When Go projects are closely related to Linux and need to use bash shell to do some initialization or query functions, Go projects will be troublesome and repetitive to write `exec.Command` . The purpose of gotestspace is to quickly create test temp directories and execute various bash commands or shell scripts in the test temp directories to help improve the efficiency of writing unit tests.	

The main entry method of the program `testspace.Create`, when this method is called, creates a shell or a workspace that the test needs to use, in which a user-defined shell is executed and the caller can use the shell directly to initialize the repository. This method supports the following options:

* WithPathOption: Used to specify the shell working directory to be executed, you need to specify an empty directory, goshellhelper will create a directory on it, and shell execution will be done in the directory
* WithTemplateOption: This is the default shell template, which provides the shell method for test_tick by default. When this template parameter is specified, a new template will be added as an append and will not override the test_tick method. It is recommended to put the public methods in this template

```shell
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}
```

**Note: test_tick: This method helps us reset the creator and committer time when creating a repository commit, ensuring that the commit number is the same and facilitating a quick git side test**

* WithEnvironmentsOption: Used to add custom environment variables, if not specified, then several environment variables will be provided by default:

```shell
GIT_AUTHOR_EMAIL=author@example.com
GIT_AUTHOR_NAME='A U Thor'
GIT_COMMITTER_EMAIL=committer@example.com
GIT_COMMITTER_NAME='C O Mitter'
```

* WithShellOption: user-defined shell to be executed, such as shell code for initializing the repository, etc.



When using Create directly, a temporary directory is created and the specified shell command is executed, and there are two ways to execute the shell in the already created workspace object:

* Execute: Execute the shell directly and return the standard output, standard errors, and error from go execution, for simple command execution

* ExecuteWithStdin: returns a pointer to a command object, which implements io.



During the creation process, gotestspace also adds a .git directory to the temporary directory to prevent dangerous git commands (such as git reset) from affecting the current project in the test area, so feel free to use it.

### How to use

* blackbox_test.go

As the naming of the file, this is a black-box test, also in the Go project, the namespace is not consistent with the current project, one more '_test', then here can only call the Go when the export method, so it is called as a black-box unit test for Go.

```go
package testspace_test

import (
	"context"
	"os"
	"path"
  ......
```

This file contains basic usage methods that can be used as a reference on how to use gotestspace

* example Directory

This directory shows how to integrate unit tests in Go, using techniques such as TestMain and SubTest