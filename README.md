# mango
Playground for an automated build tool in Go.

### Pre-requisites
 * Install Go - https://golang.org/doc/install
 * `$> go get github.com/vladimirvivien/mango`

### Using Mango
 1. Create a .build directory in your Go project
 1. Create `.build/main.go `
 1. Add `mango tasks` in your main file
 1. Build your project with `go run .build/*`

## Components
Mango uses several components that you should understand before you can use it effectively.

### The Builder
As you would imagine, this is a component responsible for `building stuff`.  The builder includes pre-defined logic (tasks) that know how to use the build system to build the stuff that it is designed for.  Today, mango provides one builder, the `Go Builder`, that can compile and test Go source code.  This builder uses the awesome GB package (https://getgb.io/) from Dave Cheney.

### Tasks
The mango `task` is how you, the person creating an automated build script, interacts with mango.  A task is basically a Go function that implements simple or complex steps that are part of the build exeuction flow.     

### The Build Files
The mango `build` files is collection of one or more Go source files.  Each file, in turn, is made up of a mix of mango tasks and other Go constructs that specify how the build flow should work.  The entry point for a build should be defined in a Go file called `main.go` with a function `func main()` defined.

### The .mango Directory
All mango build files should be stored (as a convention) in a `.mango` directory at the root level of your Go project.  The build step can then be launched using `go run .mango`.

## Using Mango
The first step is to create a directory in your Go project called `.mango`.  Next, create `.mango/main.go` with the code to build your project.  

### A Simple Build Source
The following setups the simplest way to build your Go projects.  Notice that there are no additional tasks defined. The build code adds a `GoBuilder` to the execution graph which knows how to compile Go source code.

```
package main

import "github.com/vladimirvivien/mango

func main() {
    // set the builder with defaults
    mango.Add(mango.GoBuilder)
    mango.Run()
}
```
Alternatively, you can configure the builder explicitly as shown in the following.
```
func main(){
    builder := mango.GoBuilder
    builder.Src = "./source/path"
    builder.NumCPU = 4
    builder.BuildTags = "+Something"
    mango.Add(builder)
    mango.Run()
}
```
## Mango Tasks
A task is a Go function that is executed by the mango runtime.  A task should do one thing, for instance `clean`, `test`, or `deploy`.  A task can be inserted at any point in the build graph as shown below.
```
func clean(){
    mango.RmDir("./_output")
}

func main() {
    mango.Add(clean)
    mango.Add(mango.GoBuilder)
    mango.Run()
}
```
In the previous example, the `clean` task is inserted prior to the Go builder.  This will cause the clean to execute before the Go source code is built. 
### Task Dependency
Since tasks are Go functions, their dependency can be specified by embedding function calls inside other functions.  That way, the Go compiler handles the call dependency. If a circular dependency exists, the program will naturally fail at runtime.  The following shows how this can be done.
```
var gobldr := mango.GoBuilder

func prepare() {
    clean()
    mango.Env("TEST_BUILD", "TRUE")
}

func clean() {
    mango.RmDir("./_output")
}

func main() {
    mango.Add(prepare)
    mango.Add(goBldr)
    mango.Run()
}
```
In the previous example, task `prepare` depends on `clean`.  Note, however, only `prepare` is added to the execution graph.
### Specify Tasks from CLI
By default, all tasks added to the task graph of mango will get executed.  You may, however, limit your execution to specific tasks by specifiying a task list from the command line.
```
$> go run .mango/* -t clean
```
The previous command will only execute task `clean` as defined in the mango source.
