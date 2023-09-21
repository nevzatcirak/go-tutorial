# <span style="color:#007d9c">GO Command Cheatsheet for Managing Dependencies</span>

## <span style="color:#007d9c">*Dependency Tracking*</span>

*To add your code to its **own module:***

````shell
go mod init example/module
````

This creates a ***go.mod*** file at the root of your source tree. Dependencies you add will be listed in that file.

## <span style="color:#007d9c">*Naming a Module*</span>

The location of the repository where Go tools can find **the module’s source code** (required if you’re publishing the
module). For example, it might be:

````textmate
github.com/<project-name>/
````

## <span style="color:#007d9c">*Adding a dependency*</span>

*To add **all dependencies** for a package in your module:*

````shell
go get .
````

*To add a **specific dependency:***

````shell
go get example.com/theirmodule
````

#### <span style="color:#007d9c">*Getting a specific dependency version*</span>

*To get a **specific numbered version:***

````shell
go get example.com/theirmodule@v1.2.3
````

*To get the **latest version:***

````shell
go get example.com/theirmodule@latest
````

*The following **go.mod** file **require** directive example:*

````go
require example.com/theirmodule v1.2.3
````

## <span style="color:#007d9c">*Discovering available updates*</span>

To list **all of the modules** that are **dependencies** of your current module:

````shell
go list -m -u all
````

Display the **latest version available** for a specific module:

````shell
go list -m -u example.com/theirmodule
````

## <span style="color:#007d9c">*Synchronizing your code’s dependencies*</span>

To keep your managed dependency **set tidy**, use the **go mod tidy** command.

The command has **no arguments** except for one flag, **-v**, that prints information about **removed modules.**

````shell
go mod tidy
````

## <span style="color:#007d9c">*Developing and testing against unpublished module code*</span>

You can specify that your code should use **dependency modules** that may **not be published.**

### <span style="color:#007d9c">*Requiring module code in a local directory*</span>

In the following **go.mod** file example, the current module requires the external module ***example.com/theirmodule***,
with a **nonexistent** version number (v0.0.0-unpublished) used to ensure the replacement works correctly. The replace
directive then replaces the original module path with ***../theirmodule***, a directory that is at the same level as the
current module’s directory.

```go
module example.com/mymodule

go 1.16

require example.com/theirmodule v0.0.0-unpublished

replace example.com/theirmodule v0.0.0-unpublished = >../theirmodule
```

When setting up a **require/replace pair**, use the ***go mod edit*** and ***go get*** commands to ensure that
requirements described by the file remain consistent:

````shell
go mod edit -replace=example.com/theirmodule@v0.0.0-unpublished=../theirmodule
go get example.com/theirmodule@v0.0.0-unpublished
````

### <span style="color:#007d9c">*Requiring external module code from your own repository fork*</span>

In the following **go.mod** file example:

```go
module example.com/mymodule

go 1.16

require example.com/theirmodule v1.2.3

replace example.com/theirmodule v1.2.3 = > example.com/myfork/theirmodule v1.2.3-fixed
```

````shell
$ go list -m example.com/theirmodule
example.com/theirmodule v1.2.3

$ go mod edit -replace=example.com/theirmodule@v1.2.3=example.com/myfork/theirmodule@v1.2.3-fixed
````

## <span style="color:#007d9c">*Getting a specific commit using a repository identifier*</span>

To get the module at a ***specific commit***, append the form **@commithash:**

````shell
go get example.com/theirmodule@4cf76c2
````

To get the module at a ***specific branch***, append the form **@branchname:**

````shell
go get example.com/theirmodule@bugfixes
````

## <span style="color:#007d9c">*Removing a dependency*</span>

To stop tracking all ***unused modules***, run the **go mod tidy** command.

````shell
go mod tidy
````

To ***remove*** a specific dependency, use the **go get** command

````shell
go get example.com/theirmodule@none
````

## <span style="color:#007d9c">*Specifying a module proxy server*</span>

You can set the variable to *URLs* for other module **proxy servers**, *separating URLs* with either ***a comma*** or
***a pipe***.

````shell
GOPROXY="https://proxy.example.com,https://proxy2.example.com"
````

````shell
GOPROXY="https://proxy.example.com|https://proxy2.example.com"
````

Go modules are frequently developed and distributed on version control servers and module proxies that aren’t available
on the public internet. You can set the **GOPRIVATE** environment variable to configure the go command to download and
build modules from private sources. Then the go command can download and build modules from private sources.

The **GOPRIVATE** or **GONOPROXY** environment variables may be set to lists of glob patterns matching module prefixes
that are private and should not be requested from any proxy. For example:

````shell
 GOPRIVATE=*.corp.example.com,*.research.example.com
````