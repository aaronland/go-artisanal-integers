# go-artisanal-integers

No, really.

## Caveats

This is absolutely _not_ ready for use yet. Proceed with caution.

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

	var db = flag.String("engine", "", "The name of the artisanal integer engine to use.")
	var dsn = flag.String("dsn", "", "The data source name (dsn) for connecting to the artisanal integer engine.")

	var eng artisanalinteger.Engine
	var err error

	switch *db {

	case "redis":
		eng, err = engine.NewRedisEngine(*dsn)
	case "rqlite":
		eng, err = engine.NewRqliteEngine(*dsn)
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

### Fancy

```
import (
	"flag"
	"github.com/thisisaaronland/go-artisanal-integers/util"
	"log"
)

func main() {

	var db = flag.String("engine", "", "The name of the artisanal integer engine to use.")
	var dsn = flag.String("dsn", "", "The data source name (dsn) for connecting to the artisanal integer engine.")
	var last = flag.Int("set-last-int", 0, "Set the last known integer.")
	var offset = flag.Int("set-offset", 0, "Set the offset used to mint integers.")
	var increment = flag.Int("set-increment", 0, "Set the increment used to mint integers.")
	var continuous = flag.Bool("continuous", false, "Continuously mint integers. This is mostly only useful for debugging.")

	flag.Parse()

	eng, err := util.NewArtisanalEngine(*db, *dsn)

	if err != nil {
		log.Fatal(err)
	}

	if *last != 0 {

		err = eng.SetLastInt(int64(*last))

		if err != nil {
			log.Fatal(err)
		}
	}

	if *increment != 0 {

		err = eng.SetIncrement(int64(*increment))

		if err != nil {
			log.Fatal(err)
		}
	}

	if *offset != 0 {

		err = eng.SetOffset(int64(*offset))

		if err != nil {
			log.Fatal(err)
		}
	}

	for {

		next, err := eng.NextInt()

		if err != nil {
			log.Fatal(err)
		}

		log.Println(next)

		if !*continuous {
			break
		}
	}
}	
```

## Engines

An "engine" is the interface between your code and an underlying data model (typically a database) for minting artisanal integers. The interface looks like this:

```
type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
}
```

### Memory

_Please write me_

### MySQL

_Please write me_

#### Schema

```
CREATE TABLE `integers` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `stub` char(1) NOT NULL DEFAULT '',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `stub` (`stub`)
) ENGINE=MyISAM;
```

### Redis

_Please write me_

```
$> ./bin/int -engine redis -dsn 'redis://localhost:6379'
8
```

### Rqlite

_Please write me_

This engine is not feature complete yet and still has bugs. It should be considered experimental.

```
$> ./bin/int -engine rqlite -dsn http://localhost:4001
31
```

Also, unless I am just doing it wrong when a single (following) peer is the only node left in a rqlite cluster it does not appear to know how to promote itself to be the leader. Instead it gets trapped in an endless loop like this:

```
2017/03/31 15:28:56 [WARN] raft: Election timeout reached, restarting election
2017/03/31 15:28:56 [INFO] raft: Node at 127.0.0.1:4006 [Candidate] entering Candidate state
2017/03/31 15:28:56 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4004: dial tcp 127.0.0.1:4004: getsockopt: connection refused
2017/03/31 15:28:56 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4002: dial tcp 127.0.0.1:4002: getsockopt: connection refused
2017/03/31 15:28:56 [DEBUG] raft: Votes needed: 2
2017/03/31 15:28:56 [DEBUG] raft: Vote granted from 127.0.0.1:4006. Tally: 1
2017/03/31 15:28:57 [WARN] raft: Election timeout reached, restarting election
2017/03/31 15:28:57 [INFO] raft: Node at 127.0.0.1:4006 [Candidate] entering Candidate state
2017/03/31 15:28:57 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4004: dial tcp 127.0.0.1:4004: getsockopt: connection refused
2017/03/31 15:28:57 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4002: dial tcp 127.0.0.1:4002: getsockopt: connection refused
2017/03/31 15:28:57 [DEBUG] raft: Votes needed: 2
2017/03/31 15:28:57 [DEBUG] raft: Vote granted from 127.0.0.1:4006. Tally: 1
2017/03/31 15:28:59 [WARN] raft: Election timeout reached, restarting election
2017/03/31 15:28:59 [INFO] raft: Node at 127.0.0.1:4006 [Candidate] entering Candidate state
2017/03/31 15:28:59 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4004: dial tcp 127.0.0.1:4004: getsockopt: connection refused
2017/03/31 15:28:59 [ERR] raft: Failed to make RequestVote RPC to 127.0.0.1:4002: dial tcp 127.0.0.1:4002: getsockopt: connection refused
2017/03/31 15:28:59 [DEBUG] raft: Votes needed: 2
2017/03/31 15:28:59 [DEBUG] raft: Vote granted from 127.0.0.1:4006. Tally: 1
```

But maybe I am just doing it wrong?

#### Schema

```
CREATE TABLE integers (id INTEGER PRIMARY KEY AUTOINCREMENT, stub CHAR(1))
CREATE UNIQUE INDEX `by_stub` ON integers(stub)
```

### SummitDB

Assuming the following per the [SummitDB documentation](https://github.com/tidwall/summitdb#getting-started):

```
$> ./summitdb-server
$> ./summitdb-server -p 7482 -join localhost:7481 -dir data2
$> ./summitdb-server -p 7483 -join localhost:7481 -dir data3
```

Then:

```
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7481'
2
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7481'
4
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7481'
2017/03/27 14:58:55 dial tcp 127.0.0.1:7481: getsockopt: connection refused
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7482'
2017/03/27 14:58:57 TRY 127.0.0.1:7483
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7483'
6
```

The SummitDB engine attempts to handle `TRY` messages automagically so really it looks like this:

```
$> /bin/int -engine summitdb -dsn 'redis://localhost:7482'
summitdb told me to try redis://127.0.0.1:7483 instead, so here we go...
22
```

It will also attempt to fail over to whichever peer takes over if and when the leader goes down. For example, let's say you did this and then shortly afterwards stopped the SummitDB server listening on port `7482`. You'd see something like this:

```
$> ./bin/int -engine summitdb -dsn 'redis://localhost:7482' -continuous
10282
10284
10286
10288
10290
10292
couldn't connect to leader so trying to see if the peers are rebalancing themselves (1/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (2/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (3/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (1/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (2/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (3/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (4/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (5/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (6/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (7/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (8/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (9/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (10/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (11/200)...
couldn't connect to leader so trying to see if the peers are rebalancing themselves (12/200)...
10294
10296
10298
10300
10302
10304
10306
10308
```

_Note the use of the `-continuous` flag to just keep generating integer after integer after integer..._

## Services

_Please write me_

```
type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
}
```

### Example

_Please write me_

## Tools

### int

Generate an artisanal integer on the command line.

```
./bin/int -h
Usage of ./bin/int:
  -continuous
    	Continuously mint integers. This is mostly only useful for debugging.
  -dsn string
    	The data source name (dsn) for connecting to the artisanal integer engine.
  -engine string
    	The name of the artisanal integer engine to use.
  -set-increment int
    	Set the increment used to mint integers.
  -set-last-int int
    	Set the last known integer.
  -set-offset int
    	Set the offset used to mint integers.
```

For example:

```
$> ./bin/int -engine mysql -dsn '{USER}:{PSWD}@/{DATABASE}'
182583
```

### intd

Generate an artisanal integer as a service.

```
./bin/intd -h
Usage of ./bin/intd:
  -dsn string
    	The data source name (dsn) for connecting to the artisanal integer engine.
  -engine string
    	The name of the artisanal integer engine to use.
  -host string
    	The hostname to listen for requests on (default "localhost")
  -port int
    	The port number to listen for requests on (default 8080)
  -set-increment int
    	Set the increment used to mint integers.
  -set-last-int int
    	Set the last known integer.
  -set-offset int
    	Set the offset used to mint integers.
```

For example:

```
$> ./bin/intd -engine mysql -dsn '{USER}:{PSWD}@/{DATABASE}'
```

And then:

```
$> curl localhost:8080
7001
```

## Performance

### Anecdotal

#### MySQL

Running `intd` backed by MySQL on a vanilla Vagrant machine (running Ubuntu 14.04) on a laptop against 500 concurrent users, using siege:

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

#### Redis

Running `intd` backed by Redis on a vanilla Vagrant machine (running Ubuntu 14.04) on a laptop against 1000 concurrent users, using siege:

```
siege -c 1000 http://localhost:8080
** SIEGE 3.0.5
** Preparing 1000 concurrent users for battle.
The server is now under siege...^C
Lifting the server siege...      done.

Transactions:			110761 hits
Availability:			100.00 %
Elapsed time:			63.92 secs
Data transferred:		0.59 MB
Response time:			0.06 secs
Transaction rate:		1732.81 trans/sec
Throughput:			0.01 MB/sec
Concurrency:			98.32
Successful transactions:	110761
Failed transactions:		0
Longest transaction:		6.24
Shortest transaction:		0.00
```

#### SummitDB

Running `intd` backed by SummitDB (running with [high consistency](https://github.com/tidwall/summitdb#read-consistency)) on a vanilla Vagrant machine (running Ubuntu 14.04) on a laptop against 100 concurrent users, using siege:

```
$> siege -c 100 http://localhost:8080
** SIEGE 3.0.5
** Preparing 100 concurrent users for battle.
The server is now under siege...^C
Lifting the server siege...      done.

Transactions:			418 hits
Availability:			100.00 %
Elapsed time:			44.57 secs
Data transferred:		0.01 MB
Response time:			9.13 secs
Transaction rate:		9.38 trans/sec
Throughput:			0.00 MB/sec
Concurrency:			85.65
Successful transactions:	303
Failed transactions:		0
Longest transaction:		22.80
Shortest transaction:		0.07
```

_Note: This pegged the (single) CPU on the virtual machine._

## See also

* http://www.brooklynintegers.com/
* http://www.londonintegers.com/
* http://www.neverendingbooks.org/artisanal-integers
* https://nelsonslog.wordpress.com/2012/07/29/artisinal-integers/
* https://nelsonslog.wordpress.com/2012/08/25/artisinal-integers-part-2/
* http://www.aaronland.info/weblog/2012/12/01/coffee-and-wifi/#timepixels
* https://mapzen.com/blog/mapzen-acquires-mission-integers/
* http://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/
