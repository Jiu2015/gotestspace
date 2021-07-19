## testspace

​		用来使用 go 能够快速创建 shell 执行的工作目录，以及自定义执行 shell 的工具。主要产生原因是因为以 satellite 为主的项目，需要使用 shell 来初始化测试用的仓库，但并没有太好的方法对 shell 进行统一的控制和调用，因此产生了这个辅助的项目。



程序的主入口方法 `testspace.Create` ，这个方法调用时，会创建 shell 或者测试所需要用的工作区域，在这个区域中，会执行用户自定义执行的 shell，调用者可以直接使用 shell 来初始化仓库。这个方法支持的选项：

* WithPathOption：用来指定要执行的 shell 工作目录，需要指定一个空的目录，goshellhelper 会在上面创建目录，并在目录中进行 shell 执行
* WithTemplateOption：这个是默认的 shell 模版，其中默认提供了 test_tick 的 shell 方法。当指定了这个 template 参数，则会以追加的方式添加新的 template，并不会覆盖 test_tick 方法。推荐将公共的方法放到这个 template 中

注： test_tick：这个方法在创建仓库提交的时候，可以帮助我们将创建者和提交者的时间进行重置，确保提交号一致，方便进行快速 git 方面测试

* WithEnvironmentsOption：用来添加自定义的环境变量，如果不指定，那么会默认提供几个环境变量：

```shell
GIT_AUTHOR_EMAIL=author@example.com
GIT_AUTHOR_NAME='A U Thor'
GIT_COMMITTER_EMAIL=committer@example.com
GIT_COMMITTER_NAME='C O Mitter'
```

* WithShellOption：用户自定义添加的要执行的 shell，如初始化仓库的 shell 代码等等



在使用过程中，如果不提供 options ，则所有值都有其默认值，其中工作区域中，会创建 .git 目录，用来防止在工作区中做 git 的其他操而对我们真正源码仓库产生影响，如 git reset 命令等。