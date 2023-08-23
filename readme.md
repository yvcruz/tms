# Introduction

A [Go](https://golang.org) project to establish communication with the toDus api to send messages to groups or users.

# Example

```go
package main

func main() {
	tms := NewTodusMessageService(TodusMessageServiceConfig{
		Url:      "https://broadcast.mprc.cu/api/v1/",
		Username: "test",
		Password: "test",
		Uid:      "test",
	})

	tms.SendMessageToUser("5312345678", "Hello user")
	tms.SendMessageToGroup("Hello group!!")
}
```

# Console
 > go install github.com/yvcruz/tms
 
 ```bash
 tms send --to "5312345678" --message "Hello!" --username "YOUR_USERNAME" --password "YOUR_PASSWORD"
 ```
