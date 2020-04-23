package net

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/mikelsr/bspl"
)

type protocolWrapper struct {
	// Protocol to wrap
	Protocol string `json:"protocol"`
	// Nodes the node plays
	Roles []bspl.Role `json:"roles"`
}

func wrapProtocol(p bspl.Protocol, roles ...bspl.Role) []byte {
	encoded := base64.StdEncoding.EncodeToString([]byte(p.String()))
	wrapper := protocolWrapper{Protocol: encoded, Roles: roles}
	data, _ := json.Marshal(wrapper)
	return data
}

func unwrapProtocol(data []byte) (bspl.Protocol, []bspl.Role, error) {
	var wrapper protocolWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return bspl.Protocol{}, nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(wrapper.Protocol)
	if err != nil {
		return bspl.Protocol{}, nil, err
	}
	p, err := bspl.Parse(bytes.NewReader(decoded))
	if err != nil {
		return bspl.Protocol{}, nil, err
	}
	return p, wrapper.Roles, nil
}
