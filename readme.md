# watcher

The watcher is a package that allows you to check if a keyword appears on a specific webpage. This watcher was written with caregivers in mind. When grownups are trying to keep up with sign ups, they can use this tool to check various pages without needing to look up every single page. You will need to connect to your local server http://localhost:8080/ to view the results of the search and match.

# How to import this package 

```
import "github.com/joumanae/watcher"
```

Or use `go install` through the CLI: 

```
go install github.com/joumanae/watcher/cmd/watcher@latest
```


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
