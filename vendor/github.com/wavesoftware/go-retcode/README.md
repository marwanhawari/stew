# retcode package for Go

Deterministic process exit codes based on Go errors.

## Usage

```go
err := fmt.Errorf("example error")
os.Exit(retcode.Calc(err))
```

