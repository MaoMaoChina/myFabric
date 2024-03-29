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
	"testing"

	"github.com/stretchr/testify/assert"
=======
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
>>>>>>> release-1.0
	"google.golang.org/grpc/peer"
)

type addr struct {
}

func (*addr) Network() string {
	return ""
}

func (*addr) String() string {
	return "1.2.3.4:5000"
}

func TestExtractAddress(t *testing.T) {
	ctx := context.Background()
	assert.Zero(t, ExtractRemoteAddress(ctx))

	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr: &addr{},
	})
	assert.Equal(t, "1.2.3.4:5000", ExtractRemoteAddress(ctx))
}
