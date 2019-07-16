package resolver

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type out struct {
	val bytes.Buffer
}

func (o out) String() string {
	// clear empty space
	out := func() string {
		outs := strings.Split(o.val.String(), "\n")
		length := len(outs)
		return strings.Join(outs[0:length-1], "\n")
	}()
	return out
}

func (o out) Map() map[string]interface{} {
	mapValue := make(map[string]interface{})
	rawValue := scanResult(o.val.String())
	for _, raw := range rawValue {
		key := strings.Replace(strings.Split(raw, ":")[0], `"`, "", -1)
		value := strings.Replace(strings.Split(raw, ":")[1], `"`, "", -1)
		mapValue[key] = value
	}
	return mapValue
}

func scanResult(res string) []string {
	var resp []string
	var str string
	for key, r := range strings.Replace(res, "\n", "", -1) {
		if string(r) != "{" && string(r) != "}" {
			str += string(r)
			if string(r) == "," || key == len(res)-2 {
				resp = append(resp, strings.TrimSpace(str))
				str = ""
			}
		}

		if string(r) == "}" {
			resp = append(resp, strings.TrimSpace(str))
		}
	}

	return resp
}

func Version() (string, error) {
	out, err := execute("-version")
	return out.String(), err
}

func ListServices(protoPath string) ([]string, error) {
	out, err := execute("-import-path", "./protos", "-proto", protoPath, "list")
	if err != nil {
		return nil, err
	}
	return strings.Split(out.String(), "\n"), nil
}

func Invoke(grpcServer, protoPath, service, method, data string) (map[string]interface{}, error) {
	out, err := execute("-plaintext", "-import-path", "./protos", "-proto", protoPath, "-d", data, grpcServer, service+"/"+method)
	if err != nil {
		return nil, err
	}
	return out.Map(), nil
}

func execute(args ...string) (*out, error) {
	absolutePath, _ := os.Getwd()
	splt := strings.Split(absolutePath, "/")
	path := splt[len(splt)-1]
	absoluteBinPath := fmt.Sprintf("%s/src/github.com/sog01", os.Getenv("GOPATH"))
	absoluteDirPath := strings.Join(splt[0:len(splt)-1], "/")
	if absoluteBinPath == absoluteDirPath && path != "grpc-resolver" {
		absolutePath = absoluteDirPath + "/grpc-resolver"
	} else {
		absolutePath += "/vendor/github.com/sog01/grpc-resolver"
	}

	cmd := exec.Command(fmt.Sprintf("%s/bin/grpcurl", absolutePath), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stdout
	cmd.Stdout = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stdout.String())
	}
	return &out{
		val: stderr,
	}, nil
}
