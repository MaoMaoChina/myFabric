/*
Copyright IBM Corp. All Rights Reserved.
<<<<<<< HEAD

=======
>>>>>>> release-1.0
SPDX-License-Identifier: Apache-2.0
*/

package util

import (
<<<<<<< HEAD
	"context"

=======
	"golang.org/x/net/context"
>>>>>>> release-1.0
	"google.golang.org/grpc/peer"
)

func ExtractRemoteAddress(ctx context.Context) string {
	var remoteAddress string
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if address := p.Addr; address != nil {
		remoteAddress = address.String()
	}
	return remoteAddress
}
