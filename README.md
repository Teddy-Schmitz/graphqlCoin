# GraphqlCoin

No its not another altcoin, its a graphql server to talk to a bitcoin compatible JSON-RPC server.

This is **alpha** quality software, expect bugs and the api to change. MRs are very appreciated

## Requirements

A coin daemon running with txindex turned on if you want to be able to query for random transactions.  You can set `-txindex=1` or put `txindex=1` in the conf file before starting the daemon for the first time.

## Getting Started
### Install

`go get github.com/Teddy-Schmitz/graphqlCoin`

### Docker

`docker run -p 5050:5050 -e RPCUSER=user -e RPCPASSWORD=password -e DAEMON=localhost:42068 tschmitz/graphqlCoin`

### Usage

Provide the options necessary to connect to the coin daemon, can be done with CLI, Environment variables or config file.
Config file can be in JSON,TOML,YAML,HCL,Java properties file.  Make sure to use the appropriate extension.

```
#CLI
--rpcuser=user --rpcpassword=password --daemon=localhost:42068

#ENV
RPCUSER=user
RPCPASSWORD=password 
DAEMON=localhost:42068

#CONF (called graphqlcoin(.extension) in cwd, or in $HOME/.graphqlcoin)
RPCUSER: user
RPCPASSWORD: password 
DAEMON: localhost:42068
```

###Options

All configuration options

```
REQUIRED:
rpcuser - username for jsonrpc connection to coin daemon
rpcpassword - password for jsonrpc connection to coin daemon
daemon - hostname:port to jsonrpc server

OPTIONAL:
debug - enable debug logging, default false
profiler - enable pprof on localhost:6060, default false
```

Found this useful?

```
ETH: 0xEaC8A611601F05D19816Ff502Fb0c0088ca7b9B5  
BTC: 3PnjM9BJZF64D3uCpMDQHJjrfRLeF2LM7G
BCH: 1MuD6zRHHsjfFRQjQyLsgpRKBEnA6wga81
```
