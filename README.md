# P4 golang lib

`P4 golang lib` is a awesome lib for golang to operate p4. It's not required to use `expect` script for `p4 login`.

It's a wrapper for `p4` commands, and format its output as golang structure .

```go
package main

import (
	"fmt"

	"github.com/eyotang/p4"
)


func main() {
    var (
        address  = "localhost:1666"
		user     = "tangyongqiang"
		password = "123456"
	)

	conn, err := p4.NewConn(addr, user, pass)
	if err != nil {
		fmt.Printf("New connection is failed! err: %+v", err)
	}
	result, err := conn.Dirs([]string{"//depot/*@700"})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("result: %s\n", result)
}
```

**NOTE:**  If using supervisor to start your application, please set `HOME` environment first in supervisor configure file. Don't use `~` or `"%(ENV_HOME)s"`, just use as following:

```ini
[program:demo]
environment=HOME="/root"
command=/data/app/demo/bin/demo
directory=/data/app/demo
stdout_logfile=/data/app/demo/log/stdout.log
stdout_logfile_backups=10
stderr_logfile=/data/app/demo/log/stderr.log
stderr_logfile_backups=10
```

