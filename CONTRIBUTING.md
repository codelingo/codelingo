# Contributor Guide

## Good Package Name

It's helpful if everyone using the package can use the same name 
to refer to its contents, which implies that the package name should 
be good: short, concise, evocative. By convention, packages are 
given lower case, single-word names; there should be no need for 
underscores or mixedCaps. Err on the side of brevity, since everyone 
using your package will be typing that name. And don't worry about 
collisions a priori. The package name is only the default name for 
imports; it need not be unique across all source code, and in the 
rare case of a collision the importing package can choose a different 
name to use locally. In any case, confusion is rare because the file 
name in the import determines just which package is being used.


## no Get at the start of Getters name

It's neither idiomatic nor necessary to put Get into the getter's name. If you have a field called owner (lower case, unexported), the getter method should be called Owner (upper case, exported), not GetOwner. 


## Update Comment First Word as Subject

Doc comments work best as complete sentences, which allow a wide variety of automated presentations.
The first sentence should be a summary that starts with the name being declared.


## Avoid Annotations in Comments

Comments do not need extra formatting such as banners of stars. The generated output
may not even be presented in a fixed-width font, so don't depend on spacing for alignment—godoc, 
like gofmt, takes care of that. The comments are uninterpreted plain text, so HTML and other 
annotations such as _this_ will reproduce verbatim and should not be used. One adjustment godoc 
does do is to display indented text in a fixed-width font, suitable for program snippets. 
The package comment for the fmt package uses this to good effect.


## Comment First Word as Subject

Doc comments work best as complete sentences, which allow a wide variety of automated presentations.
The first sentence should be a summary that starts with the name being declared.


## Reuse the variable name in a type switch

If the switch declares a variable in the expression, the variable will have the corresponding type in each clause. It's also idiomatic to reuse the name in such cases, in effect declaring a new variable with the same name but a different type in each case.


## Unnecessary Else

When an if statement doesn't flow into the next statement—that is, the body ends in break, continue, goto, or return—the unnecessary else is omitted.


## Defer Close File

Deferring a call to a function such as Close has two advantages. First, it guarantees that you will never forget to close the file, a mistake that's easy to make if you later edit the function to add a new return path. Second, it means that the close sits near the open, which is much clearer than placing it at the end of the function.
TODO if a function returns the open file, follow it and check if it is closed


## test

The loop variable is reused for each iteration, so it is shared across all goroutines. We need to make sure that it is unique for each goroutine. One way to do that, is passing the loop variable as an argument to the closure in the goroutine.
Note: This tenet assumes that loop variables are not shadowed inside goroutine. We need ssa to work to find the right loop variables in that case.


## Single Method Interface Name

By convention, one-method interfaces are named by the method name plus an -er suffix 
or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.

There are a number of such names and it's productive to honor them and the function names they capture. 
Read, Write, Close, Flush, String and so on have canonical signatures and meanings. To avoid confusion, 
don't give your method one of those names unless it has the same signature and meaning. Conversely, 
if your type implements a method with the same meaning as a method on a well-known type, give it the 
same name and signature; call your string-converter method String not ToString.


## Initialize instance using composite literal

Sometimes the zero value isn't good enough and an initializing constructor is necessary. We can simplify the code using a composite literal, which is an expression that creates a new instance each time it is evaluated.


## Package Comment

Every package should have a package comment, a block comment preceding the package clause. 
For multi-file packages, the package comment only needs to be present in one file, and any one will do. 
The package comment should introduce the package and provide information relevant to the package as a 
whole. It will appear first on the godoc page and should set up the detailed documentation that follows.


