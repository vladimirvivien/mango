# mango
Playground for an automated build tool in Go.

### Pre-requisites
 * Install Go - https://golang.org/doc/install
 * `$> go get github.com/vladimirvivien/mango`

### Using Mango
 1. Create a .build directory in your Go project
 1. Create `.build/main.go `
 1. Add `mango tasks` in your main file
 1. Build your project with `go run .build`

## Components
Mango uses several components that you should understand before you can use it effectively.

### The Builder
As you would imagine, this is a component responsible for `building stuff`.  The builder knows how to configure and launch its build system as implemented by the builder provider.  Today, mango provides one builder, the `Go Builder`, that can compile and test Go source code.  This builder uses the awesome GB tool (https://getgb.io/) from Dave Cheney.

### Tasks
The mango `task` is how you, the person creating an automated build script, interacts with mango.  A task is basically a Go function that implements simple or complex steps that are part of the build exeuction flow.     

### The Build Files
The mango `build` files is collection of one or more Go source files.  Each file, in turn, is made up of a mix of mango tasks and other Go constructs that specify how the build flow should work.  The entry point for a build should be defined in a Go file called `main.go` with a function `func main()` defined.

### The .build Directory
All mango build files should be stored (as a convention) in a `.build` at the root level of the project.  The build step can then be launched using `go run .build`.
