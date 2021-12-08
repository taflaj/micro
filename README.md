# micro
An implementation of microservices

## Random Numbers and Strings ##

`curl -X GET get/random/<type>[/<length>]`

Type can be `alpha`, `alphanum`, `any`, `hex`, `number`, or `special`.
Default length is 32 characters.

## Version ##

`curl -X GET get/<service>/version`

## Not Yet Implemented ##

Upload a file: `curl -X POST add/digest/file -F 'file=@</full/path/to/file>'`

Add json: `curl -X PUT add/<service>/json -H "Content-Type: application/json" -d "<json-text>"`

### Access Control ###

Get access level to a given service: `curl -X GET get/access/<service>/<ip address>[/<format>]`

Format is `text` (default) or `json`.

Response is `no` for no access, `ro` for read only, `wo` for write only, `rw` for read & write, or `default` in case access is not defined. Requestor must have read access to the access service.

Set specific access: `curl -X PUT set/access/<service>/<ip address> -d "access=<access>" [-d "owner=<owner>] [-d "remarks=<remarks>]`

Access can be `no` for no access, `ro` for read only, `wo` for write only, and `rw` for read & write. Requestor must have write access to the access service.

Set default access: `curl -X DELETE set/access/<service>/<ip address>`

Requestor must have write access to the access service.

### Message Digest ###

Download the full digest: `curl -X GET get/digest/full`

Add one or more lines: `curl -X PUT add/digest/lines[/<format>] -F 'lines[]=<text>' [-F 'lines[]=<text>' [...]]`

Format is `text` (default) or `json`.

The consumer's IP address must be registered with the system to read from and/or write to the message digest service.

### Public Keys ###

Get a public key given its email address: `curl -X GET get/pubkey/<email>`

Public keys, due to security concerns, can only be updated manually.

### Exceptions ###

A given service may need to query the access level the requestor has. In this case, it must pass a `Message` to the pubsub router, setting all fields accordingly, with attention to the following:

* `To` must include "service."
* `Service` must be set to "service."
* `IP` must be set to -1.
* `Data` must contain an `AccessLevelQuery` object.

```go
type AccessLevelQuery struct {
        Address int
        Service string
}
```

* `Address` contains the IP address of the consumer.
* `Service` contains the name of the service.

The response comes in this form:

```go
type AccessLevel struct {
	Address  int
	Service  string
	Defined  bool
	Level    string
	CanRead  bool
	CanWrite bool
}
```

* `Address` contains the IP address of the consumer.
* `Service` contains the name of the service.
* `Defined` specifies whether specific access is defined or custom (`true`) or default (`false`). If custom:
    * `Level` contains the access level: `no`, `wo`, `ro`, or `rw`.
    * `CanRead`, if `true`, indicates the requestor can read from the service.
    * `CanWrite`, if `true`, indicates the requestor can write to the service.## For Future Use ##

```
&messaging.Message{
        To:[]string{"digest"},
        CC:[]string{"logger"},
        From:"main",
        Service:"digest",
        Request:"POST",
        Command:[]string{"add", "digest", "file"},
        Data:[]string{
                "--",
                "\r\nContent-Disposition: form-data; name=\"file\";
            filename=\"Thanks\"\r\nContent-Type: application/octet-stream\r\n\r\nWe 
            would like to thank:\n\n\nVlad R. of the vdmsound project for excellent 
            sound blaster info.\nTatsuyuki Satoh of the Mame Team for making an 
            excellent FM emulator.\nJarek Burczynski for the new OPL3 emulator.\nKen 
            Silverman for his work on an OPL2 emulator.\n\nThe Bochs and DOSemu 
            projects which I used for information.\nFreeDOS for ideas in making my 
            shell.\n\nPierre-Yves G\xe9rardy for hosting the old Beta Board.\nColin 
            Snover for hosting our forum.\n\nSourceforge for hosting our homepage and 
            other development tools.\nMirek Luza, for his moderation of the forums.
            \neL_Pusher, DosFreak and MiniMax for their moderation of VOGONS forum.
            \n\ncrazyc, gulikoza, M-HT for their work on the dynrec core.\n\nJantien 
            for the version management.\nShawn, Johannes and Marcus for creating the 
            MAC OS X version.\nJochen for creating the OS/2 version.\nIdo Beeri for 
            the icon.\nripsaw8080 for his hard debugging work.\nGOG Team for the 
            splash screen.\nAll the people who submitted a bug.\nThe Beta Testers.
            \n\n\r\n--",
                "--\r\n"
        }
}
```

```
&messaging.Message{
        To:[]string{"service"},
        CC:[]string{"logger"},
        From:"main",
        Service:"service",
        Request:"PUT",
        Command:[]string{"add", "service", "fields"},
        Data:[]string{
                "--",
                "\r\nContent-Disposition: form-data; name=\"a\"\r\n\r\none field\r\n--",
                "\r\nContent-Disposition: form-data; name=\"b\"\r\n\r\nanother field\r\n--",
                "\r\nContent-Disposition: form-data; name=\"z\"\r\n\r\nand so on\r\n--",
                "--\r\n"
        }
}
```

```go
func main() {
        s1 := strings.Trim("\r\nContent-Disposition: form-data; name=\"a\"\r\n\r\none field\r\n--", "\r\n")
        fmt.Println(s1)
        s2 := strings.Split(s1, ";")
        fmt.Printf("%#v\n", s2)
        s3 := strings.Trim(s2[1], " -\r\n")
        fmt.Printf("%#v\n", s3)
        s4 := strings.Split(s3, "\r")
        fmt.Printf("%#v\n", s4)
        fmt.Printf("|%v|%v|\n", s4[0], s4[2][1:])
}
```

# Architecture #

## Entry Point ##

The entry point to all the services is by default listening on port 9999. All it does is assemble a message based on the incoming request and pass it on to the pubsub router, and then return the main response to the caller.

## PubSub Router ##

The router is listening on port 9998 by default. It receives a message through gRPC, forwards it to the subscribed services, collects the responses, and returns them as an array to the caller.

Before a service can be called, it must register itself with the router, and then unregister before it becomes inactive.

## Internal Message Format ##

This is the format that all services use to send messages to one another.

```go
type Message struct {
        To      []string
        From    string
        Service string
        Method  string
        Request string
        Command []string
        IP      uint32
        Data    string
}
```

Each message is passed as a json object using gRPC.

* `To` contains the names of all primary services that must receive this message. Each recipient is expected to respond.
* `From` contains the name of the sender of the message.
* `Service` contains the name of the service being invoked. It follows `get`, `set`, etc. on the request.
* `Method` contains the http request method.
* `Request` contains the entire http request.
* `Command` contains each component of the request (whatever is separated by `/`) in an array.
* `IP` contains the IP address of the requestor. Could be used for whitelisting and blacklisting. If set to -1, it means the request was originated by another service.
* `Data` contains whatever form data (specified using `-d` or `-F`) is included. It is encoded as JSON text.
