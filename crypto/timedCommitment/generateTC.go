package timedCommitment

import (
	"RB/config"
	"RB/crypto/binaryquadraticform"
	"RB/util"
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateTC() (*big.Int, *binaryquadraticform.BQuadraticForm, *binaryquadraticform.BQuadraticForm, *binaryquadraticform.BQuadraticForm, *binaryquadraticform.BQuadraticForm, *binaryquadraticform.BQuadraticForm, *binaryquadraticform.BQuadraticForm, *big.Int) {
	// read config file
	a, b, c := config.GetGroupParameter()
	fmt.Println(a, b, c)
	m_k_a, m_k_b, m_k_c, m_kSub_a, m_kSub_b, m_kSub_c, r_k_a, r_k_b, r_k_c := config.GetPublicGroupParameter()
	fmt.Println(m_k_a, m_k_b, m_k_c, m_kSub_a, m_kSub_b, m_kSub_c, r_k_a, r_k_b, r_k_c)
	bigTwo := big.NewInt(2)

	// get public class group
	g, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(a)), big.NewInt(int64(b)), big.NewInt(int64(c)))
	fmt.Printf("===>[GenerateTC]The group element g is (a=%v,b=%v,c=%v,d=%v)\n", g.GetA(), g.GetB(), g.GetC(), g.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	d := g.GetDiscriminant()

	m_k, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(m_k_a)), big.NewInt(int64(m_k_b)), big.NewInt(int64(m_k_c)))
	fmt.Printf("===>[GenerateTC]The group element m_k is (a=%v,b=%v,c=%v,d=%v)\n", m_k.GetA(), m_k.GetB(), m_k.GetC(), m_k.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	m_kSub, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(m_kSub_a)), big.NewInt(int64(m_kSub_b)), big.NewInt(int64(m_kSub_c)))
	fmt.Printf("===>[GenerateTC]The group element m_kSub is (a=%v,b=%v,c=%v,d=%v)\n", m_kSub.GetA(), m_kSub.GetB(), m_kSub.GetC(), m_kSub.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	r_k, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(r_k_a)), big.NewInt(int64(r_k_b)), big.NewInt(int64(r_k_c)))
	fmt.Printf("===>[GenerateTC]The group element m_k is (a=%v,b=%v,c=%v,d=%v)\n", r_k.GetA(), r_k.GetB(), r_k.GetC(), r_k.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}

	// calculate the upper bound of alpha
	nSqrt := new(big.Int)
	dAbs := new(big.Int)
	dAbs.Abs(d)
	upperBound := new(big.Int)
	nSqrt.Sqrt(dAbs)
	nSqrt.Add(nSqrt, nSqrt)
	upperBound.Sub(dAbs, nSqrt)
	upperBound.Div(upperBound, big.NewInt(10000))
	fmt.Printf("===>[GenerateTC]Upper bound of alpha is %v.\n", upperBound)

	// get random alpha
	alpha, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate alpha failed:%s", err))
	}
	fmt.Println("===>[GenerateTC]alpha is", alpha)

	// generate random message
	upper := new(big.Int)
	upper.Exp(bigTwo, big.NewInt(50), nil)
	msg, err := rand.Int(rand.Reader, upper)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTimedCommitment]Generate random message failed:%s", err))
	}
	fmt.Println("===>[GenerateTC]msg is", msg)

	// xor msg and R_k, gets c
	h, err := g.Exp(alpha)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	M_k, err := m_k.Exp(alpha)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	M_kSub, err := m_kSub.Exp(alpha)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	R_k, err := r_k.Exp(alpha)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}

	// TODO
	// F function to de modeified
	hashedRk := new(big.Int).SetBytes(util.Digest((R_k)))
	maskedMsg := new(big.Int)
	maskedMsg.Xor(msg, hashedRk)

	// evalute a series of parameters(a1, a2, a3, z) for verification
	wupperBound := new(big.Int)
	wupperBound.Mul(upperBound, dAbs)
	w, _ := rand.Int(rand.Reader, wupperBound)
	a1, err := g.Exp(w)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	a2, err := m_kSub.Exp(w)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	a3, err := m_k.Exp(w)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GenerateTC]Generate new BQuadratic Form failed: %s", err))
	}
	fmt.Println(a2.GetA(), a2.GetB(), a2.GetC())

	gHash := new(big.Int).SetBytes(util.Digest((g.GetA())))
	hHash := new(big.Int).SetBytes(util.Digest((h.GetA())))
	mkHash := new(big.Int).SetBytes(util.Digest(m_k.GetA()))
	mkSubHash := new(big.Int).SetBytes(util.Digest(m_kSub.GetA()))
	a1Hash := new(big.Int).SetBytes(util.Digest(a1.GetA()))
	a2Hash := new(big.Int).SetBytes(util.Digest(a2.GetA()))
	a3Hash := new(big.Int).SetBytes(util.Digest(a3.GetA()))

	e := big.NewInt(0)
	e.Xor(e, gHash)
	e.Xor(e, hHash)
	e.Xor(e, mkHash)
	e.Xor(e, mkSubHash)
	e.Xor(e, a1Hash)
	e.Xor(e, a2Hash)
	e.Xor(e, a3Hash)
	e.Mod(e, upperBound)

	z := new(big.Int).Set(w)
	alphaE := new(big.Int).Set(e)
	alphaE.Mul(alphaE, alpha)
	z.Sub(z, alphaE)

	fmt.Println(alpha, w, e, z)
	return maskedMsg, h, M_kSub, M_k, a1, a2, a3, z
}
