## gotestspace

​		gotestspace 是用来使用 go 能够快速创建 shell 执行的工作目录，以及自定义执行 shell 的工具，能够帮助你快速创建独立的工作区去做单元测试，提高单元测试编写效率，提升开发人员幸福感。

​	当 Go 的开发项目与 Linux 关系密切的时候，需要使用 bash shell 来做一些初始化或者查询功能的时候，Go 项目就会很麻烦，并且重复性的编写 exec.Command 以及各种环境变量等等，因此 gotestspace 工具诞生，目的是能够快速的创建测试临时目录，并在测试临时目录中，执行各种所需的 bash 命令或者 shell 脚本，帮助提升编写单元测试效率。

程序的主入口方法 `testspace.Create` ，这个方法调用时，会创建 shell 或者测试所需要用的工作区域，在这个区域中，会执行用户自定义执行的 shell，调用者可以直接使用 shell 来初始化仓库。这个方法支持的选项：

* WithPathOption：用来指定要执行的 shell 工作目录，需要指定一个空的目录，goshellhelper 会在上面创建目录，并在目录中进行 shell 执行
* WithTemplateOption：这个是默认的 shell 模版，其中默认提供了 test_tick 的 shell 方法。当指定了这个 template 参数，则会以追加的方式添加新的 template，并不会覆盖 test_tick 方法。推荐将公共的方法放到这个 template 中

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

注： test_tick：这个方法在创建仓库提交的时候，可以帮助我们将创建者和提交者的时间进行重置，确保提交号一致，方便进行快速 git 方面测试

* WithEnvironmentsOption：用来添加自定义的环境变量，如果不指定，那么会默认提供几个环境变量：

```shell
GIT_AUTHOR_EMAIL=author@example.com
GIT_AUTHOR_NAME='A U Thor'
GIT_COMMITTER_EMAIL=committer@example.com
GIT_COMMITTER_NAME='C O Mitter'
```

* WithShellOption：用户自定义添加的要执行的 shell，如初始化仓库的 shell 代码等等



在直接使用 Create 的时候，就会创建临时的目录，并执行指定的 shell 命令，在已经创建好的工作区对象中，也有两种执行 shell 的方式：

* Execute：直接执行 shell，并将标准输出，标准错误，以及 go 执行的 error 返回，适用于简单的命令执行

* ExecuteWithStdin：会返回一个 command 对象的指针，这个对象实现了 io.Reader 和 io.Writer，因此自由度会大一些，对于有持续接收标准输入的命令，则可以使用



gotestspace 在创建过程中，还会给临时目录添加 .git 目录，目的是防止在测试区域中使用危险的 git 命令（如 git reset 之类）影响当前的项目，因此请放心大胆的使用。

### 如何使用

* blackbox_test.go

正如文件的命名，这个是黑盒测试，也是 Go 项目中，命名空间与当前项目不一致，多一个 ‘_test'，那么这里只能调用 Go 当中的导出方法，因此被称作为 Go 的黑盒单元测试。

```go
package testspace_test

import (
	"context"
	"os"
	"path"
  ......
```

在这个文件中，包含了基本使用的方法，可以用来参考如何使用 gotestspace

* example 目录

这个目录中展示了如何在 Go 的单元测试进行集成，利用 TestMain 以及 SubTest 等技巧