package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mikelsr/bspl"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// ContactInfo stores the peerID, last known multiaddr
// and Services provided by a node
type ContactInfo struct {
	ID        peer.ID          `json:"peer_id"`
	MultiAddr string           `json:"multiaddr"`
	Services  []serviceWrapper `json:"services"`
}

// NewContactInfo creates a ContactInfo instance given an id, addr and services
func NewContactInfo(id peer.ID, addr multiaddr.Multiaddr, services ...Service) ContactInfo {
	data, _ := addr.MarshalJSON()
	ma := base64.StdEncoding.EncodeToString(data)
	wrappers := make([]serviceWrapper, len(services))
	for i, s := range services {
		wrappers[i] = newServiceWrapper(s)
	}
	return ContactInfo{
		ID:        id,
		MultiAddr: ma,
		Services:  wrappers,
	}
}

// Service provided by a peer
type Service struct {
	Protocol bspl.Protocol
	Roles    []bspl.Role
}

// NewService constructs a service wrapper given a protocol and roles
func NewService(protocol bspl.Protocol, roles ...bspl.Role) Service {
	return Service{
		Protocol: protocol,
		Roles:    roles,
	}
}

// serviceWrapper wraps a service provided by a node
type serviceWrapper struct {
	Protocol string      `json:"protocol"`
	Roles    []bspl.Role `json:"roles"`
}

// newServiceWrapper constructs a service wrapper given a Service
func newServiceWrapper(service Service) serviceWrapper {
	return serviceWrapper{
		Protocol: base64.StdEncoding.EncodeToString([]byte(service.Protocol.String())),
		Roles:    service.Roles,
	}
}

// WrapInfo marshals a ContactInfo to bytes
func WrapInfo(c ContactInfo) []byte {
	data, _ := json.Marshal(c)
	return data
}

// UnwrapInfo unmarshals a ContactInfo from bytes
func UnwrapInfo(data []byte) (peer.ID, multiaddr.Multiaddr, []Service, error) {
	var c ContactInfo
	if err := json.Unmarshal(data, &c); err != nil {
		return "", nil, nil, err
	}
	services := make([]Service, len(c.Services))
	for i, s := range c.Services {
		decoded, err := base64.StdEncoding.DecodeString(s.Protocol)
		if err != nil {
			return "", nil, nil, err
		}
		protocol, err := bspl.Parse(bytes.NewReader(decoded))
		if err != nil {
			return "", nil, nil, err
		}
		services[i] = NewService(protocol, s.Roles...)
	}
	decoded, err := base64.StdEncoding.DecodeString(c.MultiAddr)
	if err != nil {
		return "", nil, nil, err
	}
	var addr multiaddr.Multiaddr
	if err = addr.UnmarshalJSON(decoded); err != nil {
		return "", nil, nil, err
	}
	return c.ID, addr, services, nil
}
