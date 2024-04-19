# watcher

The watcher is a package that allows you to check if a keyword appears on a specific webpage. This watcher was written with caregivers in mind. When grownups are trying to keep up with sign ups, they can use this tool to check various pages without needing to look at every single page. 


# How to use this package 

The user should create a file nammed 'checks.txt' and write down urls and keywords following this format: 

```
kayakforkids.com level 1 
summeroffun.com 3rd grade 
```

To run the package, write the following 

```
go run cmd/watcher/main.go 

```

### When you submit your project for grading, it should have: 
+ [x] A short, meaningful module name
+ [x] A simple, logical package structure
+ [x] A README explaining briefly what the package/CLI does, how to import/install it, and a couple of examples of how to use it
+ [x] An open source licence (for example MIT)
+ [x] Passing tests with at least 90% coverage, including the CLI
+ [ ] Documentation comments for all your exported identifiers
+ [ ] Executable examples if appropriate
+ [ ] A listing on pkg.go.dev
+ [x] No commented-out code
+ [x] No unchecked errors
+ [x] No 'staticcheck' warnings
+ [ ] A Go-compatible release tag (for example v0.1.0)