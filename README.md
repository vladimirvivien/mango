# mango
Playground for an automated build tool in Go.

### Pre-requisites
 * Install Go - https://golang.org/doc/install
 * `$> go get github.com/vladimirvivien/mango`

### Using Mango
 1. Create a `.mango` directory at the root of your Go project
 1. Next, create Go file `.mango/main.go `
 1. Add `mango tasks` (functions)  in main.go
 1. Run the mango build with  `go run .mango/*`

## Components
Mango uses several components that you should understand before you can use it effectively.

### The Builder
As you would imagine, this is a component responsible for `building stuff`.  The builder includes pre-defined logic (tasks) that know how to use the build system to build the stuff that it is designed for.  Today, mango provides one builder, the `Go Builder`, that can compile and test Go source code.  This builder uses the awesome GB package (https://getgb.io/) from Dave Cheney.

### Tasks
A mango `task` is the unit of work that is used to define tasks needed to be done during the build process. 
A task is basically a Go function that implements simple (or complex steps) that are part of the build exeuction flow.

### The Build Source Files
A mango `build` is comprised of a collection of one or more Go source files inside the .mango directory.  Each file, in turn, is made up of a mix of mango tasks and other Go constructs that specify how the build flow should work.  The entry point for a mango build is defined in a Go source file called `main.go` with a function `func main()` defined.

### The .mango Directory
All mango build source files are stored (as a convention) in a `.mango` directory at the root level of your Go project.  The build step can then be launched using `go run .mango`.

## Using Mango
The first step is to create a directory at the root of your Go project called `.mango`.  Next, create `.mango/main.go` with the code to build your project.  This section walks you through some scenarios for using mango to build your Go project.

### A Simple Go Build
The following setups the simplest mango build for a given project.  The build consists of a single task that prints the word `Hello!` on standard output.
```
package main

import(
   "fmt"
   "github.com/vladimirvivien/mango"
)

func main(){
    mango.Task("hello", func(){
        fmt.Println("Hello!")
    })

    // execute tasks
    mango.Run()
}
```
By default, when `mango.Run()` is called, without parameters, it executes the tasks in the order they are decalred.  In this example, task `hello` is excuted when the mango build is run as follows:
```
go run .mango/*
```

## Mango Tasks
A task is a regular Go function that is the unit of execution for mango.  A task should do one thing and do it well,  for instance `clean`, `test`, or `deploy`, etc.
A task is declared using the `mango.Task()` function so that it can be added to the build graph.  The following example declares two tasks named `clean` and `setup`.
```
func main() {
    mango.Task("clean", func(){
        mango.RmDir("./_output")
    })

    mango.Task("setup", func(){
        mango.Env("BUILD_STAGE", "TEST")
    })

    mango.Run()
}
```
When the mango code is executed, tasks `clean` is excuted first followed by  task `setup`.  

### Task Dependency
Task dependency is specified when a task is declared using the `mango.Task()` function.  
The second parameter is used to sepecify a comma-seprated list of tasks that should be executed prior to the decalred task.
For instance in the following, task `print` depends on task `setup`.  Mango will run the parent task `setup` then execute `print`.

```
func main() {
    mango.Task("print", "setup", func(){
        name := mango.GetEnv("NAME")
        fmt.Println("Hello", name)
    })

    mango.Task("setup", func(){
        mango.Env("NAME", "MANGO")
    })

    mango.Run()
}
```
You should note that in the previous example, the task graph contains two declared tasks `print` and `setup`.  Task `setup`, however will get executed twice.  This is because, by default, task `print` will run first, causing its parent task `setup` to run.  Then task`setup` will run again because it is the second registered task in the execution graph.

This can be fixed in many ways.  However, the simplest is to specify a default task as discussed in the next section.

### Specify Default Tasks
By default, when `mango.Run()` is invoked, mango will execute all registered tasks in the order they were added. That behavior is OK as long as there is no dependencies between taks.  As we explained earlier, when there are dependencies, tasks can run more than once.  One way to fix this is to specify a default list of tasks to execute  by providing function `mango.Run` with a list of tasks to execute.

For instance, the following build source will only execute task `print` which will cause its parent task `setup` to be executed right before it. 
```
func main() {
    mango.Task("print", "setup", func(){
        name := mango.GetEnv("NAME")
        fmt.Println("Hello", name)
    })

    mango.Task("setup", func(){
        mango.Env("NAME", "MANGO")
    })

    mango.Run("print")
}
```

### Specify Default Tasks from CLI
You may specify tasks to execute from the command line by providing a list of tasks that will override the task list in the mango source code.
For instance, the following will only execute registered tasks `setup` followed by `build` by overriding the task order specified in the code.
```
$> go run .mango/* --tasks "setup"
```

## Building Go Projects
Mango comes with a Go builder that you can use to build Go projects.

```
package main

import (
    "github.com/vladimirvivien/mango"
    "github.com/vladimirvivien/mango/gobuilder"
)

func main(){
    builder := gobuilder.New()
    builder.PackageRoot = "/home/user/go"
    builder.Packages = "."
    builder.NumCPU = 4
    builder.BuildTags = "+Something"

    mango.Task("go-build", builder)

    mango.Run()
}
```
