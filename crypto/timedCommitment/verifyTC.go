package timedCommitment

import (
	"RB/config"
	"RB/crypto/binaryquadraticform"
	"RB/util"
	"fmt"
	"math/big"
)

func VerifyTC(maskedMsg *big.Int, h *binaryquadraticform.BQuadraticForm, M_kSub *binaryquadraticform.BQuadraticForm, M_k *binaryquadraticform.BQuadraticForm, a1 *binaryquadraticform.BQuadraticForm, a2 *binaryquadraticform.BQuadraticForm, a3 *binaryquadraticform.BQuadraticForm, z *big.Int) bool {
	// read config file
	a, b, c := config.GetGroupParameter()
	m_k_a, m_k_b, m_k_c, m_kSub_a, m_kSub_b, m_kSub_c, r_k_a, r_k_b, r_k_c := config.GetPublicGroupParameter()
	fmt.Println(m_k_a, m_k_b, m_k_c, m_kSub_a, m_kSub_b, m_kSub_c, r_k_a, r_k_b, r_k_c)
	// bigTwo := big.NewInt(2)

	// get public class group
	g, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(a)), big.NewInt(int64(b)), big.NewInt(int64(c)))
	fmt.Printf("===>[VerifyTC]The group element g is (a=%v,b=%v,c=%v,d=%v)\n", g.GetA(), g.GetB(), g.GetC(), g.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	d := g.GetDiscriminant()

	m_k, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(m_k_a)), big.NewInt(int64(m_k_b)), big.NewInt(int64(m_k_c)))
	fmt.Printf("===>[VerifyTC]The group element m_k is (a=%v,b=%v,c=%v,d=%v)\n", m_k.GetA(), m_k.GetB(), m_k.GetC(), m_k.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	m_kSub, err := binaryquadraticform.NewBQuadraticForm(big.NewInt(int64(m_kSub_a)), big.NewInt(int64(m_kSub_b)), big.NewInt(int64(m_kSub_c)))
	fmt.Printf("===>[VerifyTC]The group element m_kSub is (a=%v,b=%v,c=%v,d=%v)\n", m_kSub.GetA(), m_kSub.GetB(), m_kSub.GetC(), m_kSub.GetDiscriminant())
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
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

	fmt.Println(e, z)

	result1, err := g.Exp(z)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	result2, err := h.Exp(e)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	comp1, _ := result1.Composition(result2)
	fmt.Println(comp1.GetA(), comp1.GetB(), comp1.GetC())
	fmt.Println(a1.GetA(), a1.GetB(), a1.GetC())

	result3, err := m_kSub.Exp(z)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	result4, err := M_kSub.Exp(e)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	comp2, _ := result3.Composition(result4)
	fmt.Println(comp2.GetA(), comp2.GetB(), comp2.GetC())
	fmt.Println(a2.GetA(), a2.GetB(), a2.GetC())

	result5, err := m_k.Exp(z)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	result6, err := M_k.Exp(e)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from VerifyTC]Generate new BQuadratic Form failed: %s", err))
	}
	comp3, _ := result5.Composition(result6)
	fmt.Println(comp3.GetA(), comp3.GetB(), comp3.GetC())
	fmt.Println(a3.GetA(), a3.GetB(), a3.GetC())

	if !comp1.Equal(a1) {
		fmt.Println("===>[VerifyTC]test1 error")
		return false
	}
	if !comp2.Equal(a2) {
		fmt.Println("===>[VerifyTC]test2 error")
		return false
	}
	if !comp3.Equal(a3) {
		fmt.Println("===>[VerifyTC]test3 error")
		return false
	}

	return true
}
