# GRPC Resolver
This resolver act as grpc client resolver. You don't need manually create grpc client from zero !. Only using this library you can resolve all grpc handler based on *.proto files. This project is still on development phase. It's far away from production ready. Therefore, don't hesistate to put any ideas or advice in issues section to improve this project better. (using [grpcurl](https://github.com/fullstorydev/grpcurl) for core engine) 

# Examples
```
package main

import (
	"fmt"
	"log"

	r "github.com/sog01/grpc-resolver"
)

func main() {
	grpcResolver, err := r.New(r.Conf{
		GrpcServer: "localhost:50051",
		ProtoPath:  "helloworld.proto",
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := grpcResolver.Exec("SayHello", "helloworld.Greeter", r.Adapter(`{"name": "juliardi"}`))
	if err != nil {
		log.Fatal("fatal: ", err)
	}
	var message string
	message = result.GetString("message")

	fmt.Println(message)
}
```