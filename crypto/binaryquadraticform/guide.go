package binaryquadraticform

import (
	"fmt"
	"math/big"
)

// Initialization
func TestInit() {
	form1, _ := NewBQuadraticForm(big.NewInt(1), big.NewInt(1), big.NewInt(6))
	fmt.Println(form1.a, form1.discriminant)

	form2, _ := NewBQuadraticFormByDiscriminant(big.NewInt(1), big.NewInt(1), big.NewInt(-100000))
	fmt.Println(form2.a, form2.discriminant)

	got, _ := form2.Composition(form2)
	fmt.Println(got.a, got.discriminant)
}

// Exp
func TestExp() {
	form1, _ := NewBQuadraticForm(big.NewInt(31), big.NewInt(24), big.NewInt(15951))

	got, _ := form1.Exp(big.NewInt(200))
	fmt.Println(got.a, got.discriminant)
}
