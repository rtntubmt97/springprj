**Start app**
On Linux, cd into the source directory, start app by typing:
```console
go run app/master/master.go
```
**Turn off logger**
Turn off logger by comment it in utils/log.go:
```go
func LogE(msg string) {
    innerLog.Output(2, fmt.Sprintf("[logE] %s\n", msg))
}
```
to
```go
func LogE(msg string) {
    // innerLog.Output(2, fmt.Sprintf("[logE] %s\n", msg))
}
```
... similar to LogI, LogD  
**Don't turn off logR**
This logger is being used to print the result, so don't turn it off
**Input**
In current version, input is read from *input.ini* file. If you want to print the input in the console, uncomment *Print* statement in app/master/master.go
```go
inputRaw, err := reader.ReadString('\n')
// fmt.Print(inputRaw)
```
to
```go
inputRaw, err := reader.ReadString('\n')
fmt.Print(inputRaw)
```
	
