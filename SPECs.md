# SPECs



### /fetch

-  Request `nil`
-  Response

```go
type Response struct {
    Serv map[Service]Status
    Noti []Notification
}
```

### /callback

- Request

```go
type Request struct {
    Heading string
    Content string
}
```

- Response `nil`

### /register

- Request

```go
type Request struct {
    Name string
    Addr string
}
```

- Response `nil`

### /revoke

- Request `Name string`
- Response `nil`



### Assistive

```go
type Service struct {
	Name string
	Addr string
}

type Status struct {
	Tm         time.Time
	Stat       int
	Informator string
}

type Notification struct {
	Tm      time.Time
	Heading string
	Content string
}
```

