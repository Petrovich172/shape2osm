package main

import (
	"fmt"
)

// Point - dimensions
type Point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

// TDsize - alias to Point
type TDsize = Point

type Tensor struct {
	Data []float64
	Size TDsize
}

// NewTensor - Constructor for Tensor type.
/*
	x - number of columns (width);
	y - number of rows (height);
*/

func NewTensor(x, y int) Tensor {
	return Tensor{
		Data: make([]float64, x*y),
		Size: TDsize{
			X: x,
			Y: y,
		},
	}
}

func (t1 *Tensor) Get(x, y int) float64 {
	return (*t1).Data[y*(*t1).Size.X+x]
}

func (t1 *Tensor) Set(x, y int, val float64) {
	(*t1).Data[y*(*t1).Size.X+x] = val
}

// SetData - Set data for *Tensor
/*
	r - number of rows;
	c - number of columns (width);
	data - 1-D array of float64.
*/
func (t1 *Tensor) SetData(c, r int, data []float64) {
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
				(*t1).Set(i, j, data[j*c+i])
				//fmt.Println(i, j, data[j*c+i])
		}
	}
}

// Print - Pretty print for *Tensor
func (t1 *Tensor) Print() {
	mx := (*t1).Size.X
	my := (*t1).Size.Y
		for y := 0; y < my; y++ {
			for x := 0; x < mx; x++ {
				fmt.Printf("%.8f\t", (*t1).Get(x, y))
			}
			fmt.Println()
		}
}


// Conv2D - apply convolution to t1 using t2 kernel
func (t1 *Tensor) Conv2D (t2 Tensor, stride [2]int, padding [2]int) Tensor {

outputX := ((*t1).Size.X - t2.Size.X + 2*padding[1])/stride[1] + 1
outputY := ((*t1).Size.Y - t2.Size.Y + 2*padding[0])/stride[0] + 1
// fmt.Println(outputX, outputY)
outputData := NewTensor(outputX, outputY)


// index for input element in output array
el := 0
iLimit := outputY
jLimit := outputX
if stride[0] > 1 {
	iLimit = outputY + stride[0] + 1
}
if stride[1] > 1 {
	jLimit = outputX + stride[1] + 1
}


	for i := 0; i < iLimit; i = (i + stride[0]) {				// rows
		for j := 0; j < jLimit; j = (j + stride[1]) {			// columns
			if i == 0 && j == 0 {
				el = 0
			} else {
				el += 1
			}
			for m := 0; m < t2.Size.X; m++ {		// kernel rows
				
				for n := 0; n <	t2.Size.Y; n++ {	// kernel columns

					// index of input signal, used for checking boundary
					ii := i + m;
					jj := j + n;

	                // ignore input samples which are out of bound
					if ii >= 0 && ii < (*t1).Size.Y && jj >= 0 && jj < (*t1).Size.X {
						outputElement := 0.0
						kernelElement := t2.Data[(m * t2.Size.X + n)]
						if i == 0 {
							inputElement := (*t1).Data[( (m) * t1.Size.X + n) + j]
							outputElement = inputElement * kernelElement
						} else {
							inputElement := (*t1).Data[ii * (*t1).Size.X + jj]
							outputElement =  inputElement * kernelElement
						}
						// fmt.Println("i,j", i, j, "\tii,jj: ", ii, jj, "\t||", (m * t2.Size.X + n), "||\t", (*t1).Data[ii * (*t1).Size.X + jj], "*", t2.Data[(m * t2.Size.X + n)])

						// Filling output array
						outputData.Data[el] += outputElement
					}
				}
			}
		}
	}
	fmt.Print("Input Data: ", t1, "\n", "Kernel: ", t2, "\n", "Output Data: ", outputData)
	return outputData
}


func main() {
	inputData := NewTensor(8, 9)
	inputData.SetData(8, 9, []float64{-0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, -0.9, -0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, -0.17, 0.18, -0.19, 0.20, 0.21, 0.22, 0.23, 0.24, -0.25, 0.26, 0.27, -0.28, 0.29, 0.30, 0.31, 0.32, -0.33, 0.34, 0.35, 0.36, -0.37, 0.38, 0.39, 0.40, -0.41, 0.42, 0.43, 0.44, 0.45, -0.46, 0.47, 0.48, -0.49, 0.50, 0.51, 0.52, 0.53, 0.54, -0.55, 0.56, -0.57, 0.58, 0.59, 0.60, 0.61, 0.62, 0.63, -0.64, -0.65, 0.66, 0.67, 0.68, 0.69, 0.70, 0.71, 0.72})

	kernel := NewTensor(3, 3)
	kernel.SetData(3, 3, []float64{0.10466029, -0.06228581, -0.43436298, 0.44050909, -0.07536250, -0.34348075, 0.16456005, 0.18682307, -0.40303048})

	stride := [2]int{4, 1}
	padding := [2]int{0, 0}
	res := inputData.Conv2D(kernel, stride, padding)
	fmt.Print("\n\n")
	res.Print()
}