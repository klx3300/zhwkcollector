# SPECs

### /fetch

-  Request

```go
type Request struct {
    Name string
    StatOnly string // "true"/"false", unidentified will be treated "false"
}
```

-  Response

```go
type Response struct {
    Serv []Service
    Stat []Status
    Noti []Notification
}
```

### /callback

- Request

```go
type Request struct {
    Heading string
    Content string
    DestNames []string
    AnyAck int // beside DestNames, we may need how much acknowledges?
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



### Encrypted

```go
type Encrypted struct {
    EncryptionType string
    DataBase64 string
}
```

For current version, only "AES" in `EncryptionType` field is effective. For bare connections, just pass bare structure.