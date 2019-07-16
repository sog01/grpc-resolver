package resolver

import (
	"errors"
)

type (
	Resolver struct {
		grpcServer string
		protoPath  string
	}

	Conf struct {
		GrpcServer string
		ProtoPath  string
	}

	Response struct {
		val   map[string]interface{}
		incre int
	}
)

func New(conf Conf) (*Resolver, error) {
	if conf.GrpcServer == "" {
		return nil, errors.New("grpc server empty")
	}
	if conf.ProtoPath != "" {
		_, err := ListServices(conf.ProtoPath)
		if err != nil {
			return nil, err
		}
	}
	return &Resolver{
		grpcServer: conf.GrpcServer,
		protoPath:  conf.ProtoPath,
	}, nil
}

func (res *Response) GetString(key string) string {
	result, _ := res.Get(key).(string)
	return result
}

func (res *Response) GetInt(key string) int {
	result, _ := res.Get(key).(int)
	return result
}

func (res *Response) Get(key string) interface{} {
	return res.val[key]
}

func (r Resolver) Exec(method, service string, req Request) (*Response, error) {
	result, err := Invoke(r.grpcServer, r.protoPath, service, method, req.Extract())
	return &Response{val: result}, err
}
