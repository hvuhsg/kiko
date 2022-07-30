package vector

import (
	"math"
	"math/rand"
)

type Vector interface {
	Add(Vector) Vector
	Sub(Vector) Vector
	Div(float64) Vector
	Mul(float64) Vector
	Dot(Vector) float64
	Norm() float64
	Norm2() float64
	Normalize() Vector
	Dimension(int) float64
	Dimensions() int
}

type vector struct {
	arr []float64
}

func NewVectorFromArray(arr []float64) Vector {
	vec := new(vector)
	vec.arr = arr

	return vec
}

func NewRandomVector(dimensions int) Vector {
	arr := make([]float64, dimensions)

	for i := range arr {
		arr[i] = rand.Float64() * 1000
	}

	return NewVectorFromArray(arr)
}

func NewZerosVector(dimensions int) Vector {
	arr := make([]float64, dimensions)

	return NewVectorFromArray(arr)
}

func (v *vector) Add(other Vector) Vector {
	addVec := new(vector)
	addVec.arr = make([]float64, len(v.arr))

	for index, value := range v.arr {
		addVec.arr[index] = value + other.Dimension(index)
	}

	return addVec
}

func (v *vector) Sub(other Vector) Vector {
	subVec := new(vector)
	subVec.arr = make([]float64, len(v.arr))

	for index, value := range v.arr {
		subVec.arr[index] = value - other.Dimension(index)
	}

	return subVec
}

func (v *vector) Div(scalar float64) Vector {
	divVec := new(vector)
	divVec.arr = make([]float64, len(v.arr))

	for index, value := range v.arr {
		divVec.arr[index] = value / scalar
	}

	return divVec
}

func (v *vector) Mul(scalar float64) Vector {
	mulVec := new(vector)
	mulVec.arr = make([]float64, len(v.arr))

	for index, value := range v.arr {
		mulVec.arr[index] = value * scalar
	}

	return mulVec
}

func (v *vector) Dot(other Vector) float64 {
	var sum float64

	for index, value := range v.arr {
		sum += value * other.Dimension(index)
	}

	return sum
}

func (v *vector) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *vector) Norm2() float64 {
	return v.Dot(v)
}

// Normalize returns a unit vector in the same direction as v.
func (v *vector) Normalize() Vector {
	n2 := v.Norm2()
	if n2 == 0 {
		return NewZerosVector(v.Dimensions())
	}
	return v.Mul(1 / math.Sqrt(n2))
}

func (v *vector) Dimension(index int) float64 {
	return v.arr[index]
}

func (v *vector) Dimensions() int {
	return len(v.arr)
}
