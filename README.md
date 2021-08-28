# grpc-streaming
Tips and tricks of authentication grpc stream server

## Inspect grpc package and services
Use evans tool for inspecting your gRPC service

    sonu@dev:evans -r -p 50005
    sonu@dev:~/go/src/github.com/grpc-server-streaming$ evans -r -p 50005
    
      ______
     |  ____|
     | |__    __   __   __ _   _ __    ___
     |  __|   \ \ / /  / _. | | '_ \  / __|
     | |____   \ V /  | (_| | | | | | \__ \
     |______|   \_/    \__,_| |_| |_| |___/
    
     more expressive universal gRPC client
    
    127.0.0.1:50005> show packages
    +-------------------------+
    |         PACKAGE         |
    +-------------------------+
    | authapi                 |
    | grpc.reflection.v1alpha |
    | sapi                    |
    +-------------------------+
    
    127.0.0.1:50005> package authapi
    
    authapi@127.0.0.1:50005> show service
    +-------------+-------+--------------+---------------+
    |   SERVICE   |  RPC  | REQUEST TYPE | RESPONSE TYPE |
    +-------------+-------+--------------+---------------+
    | AuthService | Login | LoginRequest | LoginResponse |
    +-------------+-------+--------------+---------------+
    
    authapi@127.0.0.1:50005> service AuthService
    
    authapi.AuthService@127.0.0.1:50005> call Login
    username (TYPE_STRING) => admin
    password (TYPE_STRING) => password
    {
      "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzAxODEzMjQsInVzZXJuYW1lIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4ifQ.w1Bah078OUKgdx0NgbLRuGM8nm2fImsmFRtW_lsdbVw"
    }
    
    authapi.AuthService@127.0.0.1:50005> 
    
    authapi.AuthService@127.0.0.1:50005> package sapi
    
    sapi@127.0.0.1:50005> show service
    +---------------+---------------+--------------+---------------+
    |    SERVICE    |      RPC      | REQUEST TYPE | RESPONSE TYPE |
    +---------------+---------------+--------------+---------------+
    | StreamService | FetchResponse | Request      | Response      |
    +---------------+---------------+--------------+---------------+
    
    sapi@127.0.0.1:50005> service StreamService
    
    sapi.StreamService@127.0.0.1:50005> header authorization="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzAxODEzMjQsInVzZXJuYW1lIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4ifQ.w1Bah078OUKgdx0NgbLRuGM8nm2fImsmFRtW_lsdbVw"
    
    sapi.StreamService@127.0.0.1:50005> call FetchResponse
    id (TYPE_INT32) => 1
    {
      "result": "Request #0 for Id:1"
    }
    {
      "result": "Request #1 for Id:1"
    }
    {
      "result": "Request #2 for Id:1"
    }
    {
      "result": "Request #3 for Id:1"
    }
    {
      "result": "Request #4 for Id:1"
    }
    
    127.0.0.1:50005> package sapi
    
    sapi@127.0.0.1:50005> service StreamService
    
    sapi.StreamService@127.0.0.1:50005> call FetchResponse
    id (TYPE_INT32) => 1
    command call: rpc error: code = Unauthenticated desc = Unauthorized, token not provided
    
    sapi.StreamService@127.0.0.1:50005> header authorization="EMPTY"
    
    sapi.StreamService@127.0.0.1:50005> call FetchResponse
    id (TYPE_INT32) => 1
    command call: rpc error: code = Unauthenticated desc = Unauthorized, invalid token






