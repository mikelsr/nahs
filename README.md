![NaHS logo][logo]

Network of Autonomous and Heterogeneous Services (NaHS)

[![Build Status](https://travis-ci.com/mikelsr/nahs.svg?token=736yMuj6XUy7yCEvSpBB&branch=master)](https://travis-ci.com/mikelsr/nahs)
[![codecov](https://codecov.io/gh/mikelsr/nahs/branch/master/graph/badge.svg?token=PSTZ46XN7Q)](https://codecov.io/gh/mikelsr/nahs)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mikelsr/nahs)](https://github.com/mikelsr/nahs/blob/master/go.mod)
[![GoDoc Reference](https://godoc.org/github.com/mikelsr/nahs?status.svg)](https://godoc.org/github.com/mikelsr/nahs)

## Modules

* `events`: Describes BSPL instance events according to the [implementation](https://github.com/mikelsr/bspl/tree/master/implementation). As of now there are three events:

  * `NewEvent` to create an [instance](https://github.com/mikelsr/bspl/blob/master/bspl.go#L27) of a [protocol](https://github.com/mikelsr/bspl/blob/master/bspl.go#L20).

  * `UpdateEvent` to update an instace comparing it to a future version of it.

  * `DropEvent` to cancel an instance for any reason.

* `net`: Networking components. The main struct is [`Node`](https://github.com/mikelsr/nahs/blob/master/net/node.go). A node has a [BSPL reasoner](https://github.com/mikelsr/bspl/blob/master/bspl.go#L25) and a [LibP2P host](https://github.com/libp2p/go-libp2p-core/blob/master/host/host.go), implementing methods and handlers to send BSPL components between network peers. Nodes discover each other either manually or with the libp2p implementation of rendezvous (**preferred**) using the default bootstrap nodes.

## Other folders

* `config`: Contains the private key of the main network (which is public, private only limits interaction
to NaHS nodes).

* `scripts`: Contains a script to generate a private network key.

* `test`: Test resources.

[logo]: .res/img/nahs.png "NaHS logo"
