// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"math"

	"github.com/gonum/blas/blas64"
)

// Dlarfg generates an elementary reflector for a Householder matrix. It creates
// a real elementary reflector of order n such that
//  H * (alpha) = (beta)
//      (    x)   (   0)
//  H^T * H = I
// H is represented in the form
//  H = 1 - tau * (1; v) * (1 v^T)
// where tau is a real scalar.
//
// On entry, x contains the vector x, on exit it contains v.
func (impl Implementation) Dlarfg(n int, alpha float64, x []float64, incX int) (beta, tau float64) {
	if n < 0 {
		panic(nLT0)
	}
	if n <= 1 {
		return alpha, 0
	}
	checkVector(n-1, x, incX)
	bi := blas64.Implementation()
	xnorm := bi.Dnrm2(n-1, x, incX)
	if xnorm == 0 {
		return alpha, 0
	}
	beta = -math.Copysign(impl.Dlapy2(alpha, xnorm), alpha)
	safmin := dlamchS / dlamchE
	knt := 0
	if math.Abs(beta) < safmin {
		// xnorm and beta may be innacurate, scale x and recompute.
		rsafmn := 1 / safmin
		for {
			knt++
			bi.Dscal(n-1, rsafmn, x, incX)
			beta *= rsafmn
			alpha *= rsafmn
			if math.Abs(beta) >= safmin {
				break
			}
		}
		xnorm = bi.Dnrm2(n-1, x, incX)
		beta = -math.Copysign(impl.Dlapy2(alpha, xnorm), alpha)
	}
	tau = (beta - alpha) / beta
	bi.Dscal(n-1, 1/(alpha-beta), x, incX)
	for j := 0; j < knt; j++ {
		beta *= safmin
	}
	return beta, tau
}
