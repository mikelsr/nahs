package net

import (
	"bufio"
	"context"

	"github.com/mikelsr/bspl"
	"github.com/mikelsr/nahs/events"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// Node represents a single NaHS peer.
type Node struct {
	// libp2p Host
	host host.Host
	// Contacts of the Node
	Contacts Contacts
	// context of the node and the host
	context context.Context
	// host cancelation function
	cancel context.CancelFunc
	// dht table with information about network peers
	dht *dht.IpfsDHT
	// OpenInstances maps instance keys to peer.IDs to
	// verify that the node sending the event is the one
	// who created it
	OpenInstances map[string]peer.ID
	// routing // TODO: rendezvous point?
	routing *discovery.RoutingDiscovery
	// protocols this node offers
	protocols []bspl.Protocol
	// resoner to handle BSPL logic
	reasoner bspl.Reasoner
	// roles this node plays for each protocol mapped to
	// protocol keys
	roles map[string][]bspl.Role
}

// NewNode is the default constructor for Node.
func NewNode(reasoner bspl.Reasoner, options ...libp2p.Option) *Node {
	n := newNode(options...)
	n.reasoner = reasoner
	// Connect the node to the bootstrap nodes to discover other peers
	n.configDiscovery()
	return n
}

// newNode is a constructor that requires no bspl.Reasoner
// and doesn't connect to the boostrap nodes used only inside
// this package.
func newNode(options ...libp2p.Option) *Node {
	n := new(Node)

	n.Contacts = make(Contacts)
	n.OpenInstances = make(map[string]peer.ID)
	n.roles = make(map[string][]bspl.Role)

	n.context, n.cancel = context.WithCancel(context.Background())
	// Contatenate options parameter to default options
	opt := append(options, []libp2p.Option{
		libp2p.ListenAddrStrings(listenAddrs...),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let this host use relays and advertise itself on relays
		// libp2p.EnableAutoRelay(),
		// Attempt to open ports using uPNP for NATed hosts
		libp2p.NATPortMap(),
	}...)

	// opt = append(opt, libp2p.PrivateNetwork(loadPrivNetPSK()))

	h, err := libp2p.New(n.context, opt...)
	if err != nil {
		panic(err)
	}
	n.host = h

	// set stream handlers
	n.setStreamHandlers()

	logger.Infof("Created node with ID '%s'.", h.ID())
	return n
}

// NodeFromPrivKey is a newNode wrapper to create a new Node with the specified
// private key. Additional options may be provided.
func NodeFromPrivKey(reasoner bspl.Reasoner, sk crypto.PrivKey, options ...libp2p.Option) *Node {
	n := nodeFromPrivKey(sk, options...)
	n.reasoner = reasoner
	n.configDiscovery()
	return n
}

// nodeFromPrivKey is the same as NodeFromPrivKey but requires
// no bspl.Reasoner and doesn't connect to the boostrap nodes
// used only inside this package.
func nodeFromPrivKey(sk crypto.PrivKey, options ...libp2p.Option) *Node {
	return newNode(append(options, libp2p.Identity(sk))...)
}

// ID of the libp2p host of the Node
func (n Node) ID() peer.ID {
	return n.host.ID()
}

// Addrs returns the multiaddr of the libp2p host of the Node
func (n Node) Addrs() []multiaddr.Multiaddr {
	return n.host.Addrs()
}

// AddProtocol adds a protocol to the node and establishes what roles
// the node plays in that protocol. If the protocol was already added,
// the roles that weren't already established are added.
func (n *Node) AddProtocol(p bspl.Protocol, roles ...bspl.Role) {
	playedRoles, found := n.roles[p.Key()]
	if !found {
		n.protocols = append(n.protocols, p)
		n.roles[p.Key()] = roles
		return
	}
	for _, role := range roles {
		found := false
		for _, playedRole := range playedRoles {
			if playedRole == role {
				found = true
				break
			}
		}
		if !found {
			n.roles[p.Key()] = append(n.roles[p.Key()], role)
		}
	}
}

// ExportKey returns the marshaled private key of the Node
func (n *Node) ExportKey() []byte {
	prv := n.host.Network().Peerstore().PrivKey(n.host.ID())
	b, _ := crypto.MarshalPrivateKey(prv)
	return b
}

// Reasoner returns the reasoner of the Node
func (n *Node) Reasoner() bspl.Reasoner {
	return n.Reasoner()
}

// SendEvent sends an events.Event to the target node.
// If the node is unreachable, the address is not known
// or some error occurs the error is returned. If the
// event was invalid, false is returned. If it was registered
// correctly, true is returned.
func (n *Node) SendEvent(target peer.ID, event events.Event) (bool, error) {
	data, err := event.Marshal()
	if err != nil {
		return false, err
	}
	stream, err := n.host.NewStream(n.context, target, protocolEventID)
	if err != nil {
		return false, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	rw.Write(data)
	rw.WriteByte(exchangeEnd)
	if err := rw.Flush(); err != nil {
		return false, err
	}
	return readEventResponse(rw)
}
