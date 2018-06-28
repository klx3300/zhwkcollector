# zhwkcollector
Collect something from somewhere to elsewhere..?



### Configuration File

JSON-encoded `map<string,string>`

Must-have sections are:

- NotificationSave
- ServiceSave
- ListenPort
- Timeout (in milliseconds, number written in string)
- FetchSymKey (at most 32 bytes, will trunc if larger)
- CallbackSymKey
- ControlSymKey
- ServicePollSymKey
- AESIV (must equal to the blocksize of AES. will be defaulted to repeating "IVECTOR" if not specified)
- AcceptBareConn ("true" or "false", unmatched will be treated as "false")



### Run

```shell
$ ./zhwkcollector <config-file>
```



### Change Logs

- 0.1 Basic functions completed. Not recommended to use at production environments.

