# golox: a lox interpreter written in Go
 
This is a working interpreter of the Lox language, as
used in the book Crafting Interpreters by Robert Nystrom

This interpreter is feature complete through Chapter 10
of Crafting Interpreters, meaning it implements variables, 
printing, loops, control flow, and functions (complete with
working return statements and control flow). Does NOT include variable resolution, classes, or inheritence <sub>(yuck)</sub>

# Instructions

## Installation
In order to compile the source code, you will need [Go 1.20+](https://go.dev/dl/), as well as [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git). Go installation
instructions can be found [here](https://go.dev/doc/install). To make sure the compiler has been properly installed, type the following command:
```
$ go version
```

Finally, clone the repository into the current directory with
```
$ git clone https://github.com/reilandeubank/golox
```
or into ```path/to/directory``` using 
```
$ git clone https://github.com/reilandeubank/golox path/to/directory
```

## Use
First, move into the ```golox``` directory that was just created

From here, you can compile an executable ```./main``` in ```/cmd```. 
```
go build cmd/main.go
```

Usage for the interpreter is
```
$ ./main
```
to start a REPL or
```
$ ./main file.lox
```
to run ```file.lox```
