# MPS
[Mathematical Programming System (MPS)](http://en.wikipedia.org/wiki/MPS_%28format%29) format parser for Golang.
This package reads a plain/gzip mps file and returns the [linear programming problem](http://en.wikipedia.org/wiki/Linear_programming) in [standard form](http://en.wikipedia.org/wiki/Linear_programming#Standard_form) represented as the following matrices :
1. A, constraints coefficient matrix.
2. b, column vector containing the RHS of constraints.
3. c, column vector containing objective function coefficients.

It also returns the functions that map from variables used in the standard form to the the actual variables of the problem.

This package depends on [go.matrix](https://github.com/skelterjohn/go.matrix)

[![GoDoc](http://godoc.org/github.com/dennisfrancis/mps?status.png)](http://godoc.org/github.com/dennisfrancis/mps)