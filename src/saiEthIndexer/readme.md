## Project layout
1. cmd/app - entry point of application
2. config - contains whole config and handling specific config options (common options handles in internal/config-internal). Also contains contracts.json, where contracts to control specified
3. handlers -  defined handlers 
4. internal - main framework folder
    app - main application functionality (registering config,storage,handlers and etc)
    config-internal - common config options (server settings and etx)
    http - boilerplate code for http server (routing options, middlewares)
6. tasks - main busyness logic
7. pkg - common code (start http server, eth client and etc)
8. utils - code to deal with another sai services and common utils


## config/config.json (application configuration)
- common(http_server,socket_server, web_socket) - common server options for http,socket and web socket servers
- geth-server - geth-server address
- storage - options for saiStorage
- start_block - number of block to start parsing 
- operations - commands under special control
- sleep - duration after which we get next block from geth server

## config/contracts.json (stored on control contracts)
- address - address of contract
- abi - Application Binary Interface (ABI) of a smart contract 
- start_block - number of block, from which contract is valid


## Add contract to control list command
curl -X POST <host:port>/v1/add_contract  -H "Content-Type: application/json" -d '{"contracts": [{"address": "0x9fe3Ace9629468AB8858660f765d329273D94D6D","abi": "324234","start_block":123},{"address": "0x9fe3Ace9629468AB8858660f765d329273D94D6E","abi":"test","start_block":34}]}'

## Delete contract from control list command
curl -X POST <host:port>/v1/delete_contract  -H "Content-Type: application/json" -d '{"addresses": ["0x9fe3Ace9629468AB8858660f765d329273D94D6E","0x9fe3Ace9629468AB8858660f765d329273D94D6W"]}'
