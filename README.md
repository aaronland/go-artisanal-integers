# go-artisanal-integers

No, really.

## Caveats

This is absolutely _not_ ready for use yet. For example interfaces (for integer `Engine` thingies) have not been finalized either.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

### Simple

```
package main

import (
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
	"log"
)

func main() {

	var db = flag.String("engine", "", "...")
	var dsn = flag.String("dsn", "", "...")

	var eng artisanalinteger.Engine
	var err error

	switch *db {

	case "redis":
		eng, err = engine.NewRedisEngine(*dsn)
	case "summitdb":
		eng, err = engine.NewSummitDBEngine(*dsn)
	case "mysql":
		eng, err = engine.NewMySQLEngine(*dsn)
	default:
		log.Fatal("Invalid engine")
	}

	if err != nil {
		log.Fatal(err)
	}

	next, err := eng.NextId()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(next)
}
```

## Engines

_Please write me_

```
type Engine interface {
	NextId() (int64, error)
	LastId() (int64, error)
	SetLastId(int64) error
	SetOffset(int64) error
	SetIncrement(int64) error
}
```

### MySQL

### Redis

### SummitDB

## Services

_Please write me_

```
type Service interface {
	NextId() (int64, error)
	MaxId() (int64, error)
}
```

### Example

## Tools

### int

Generate an artisanal integer on the command line.

```
$> ./bin/int -engine mysql -dsn '{USER}:{PSWD}/{DATABASE}'
182583
```

### intd

Generate an artisanal integer as a service.

```
$> ./bin/intd -engine mysql -dsn '{USER}:{PSWD}/{DATABASE}'
```

And then

```
$> curl localhost:8080
7001
```

## Performance

### Anecdotal

Running `intd` on a vanilla Vagrant machine (running Ubuntu 14.04) on a laptop against 500 concurrent users, using siege:

```
$> siege -c 500 http://localhost:8080
** SIEGE 3.0.5
** Preparing 500 concurrent users for battle.
The server is now under siege...^C
Lifting the server siege...      done.

Transactions:			58285 hits
Availability:			100.00 %
Elapsed time:			70.71 secs
Data transferred:		0.32 MB
Response time:			0.02 secs
Transaction rate:		824.28 trans/sec
Throughput:			0.00 MB/sec
Concurrency:			14.98
Successful transactions:	58217
Failed transactions:		0
Longest transaction:		1.70
Shortest transaction:		0.00
```

## See also

