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
	(*t1).Data[y*(*t1).Size.X+x] = val //Set(i, j, data[j*c+i])
	//fmt.Println((*t1).Data[y*(*t1).Size.X+x])
	//fmt.Println(y,x)
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

outputX := ((*t1).Size.X - t2.Size.X + 2*padding[0])/stride[0] + 1
outputY := ((*t1).Size.Y - t2.Size.Y + 2*padding[1])/stride[1] + 1
// fmt.Println(outputX, outputY)
outputData := NewTensor(outputX, outputY)

// find center position of kernel (half of kernel size)
kCenterX := t2.Size.X / 2;
kCenterY := t2.Size.Y / 2;

// index to input element into output array
el := 0


	for i := 0; i < outputX; i = (i + stride[0]) {              // rows
	    for j := 0; j < outputY; j = (j + stride[1]) {          // columns
	    	if i == 0 && j == 0 {
	    		el = 0
	    	} else {
	    		el += 1
	    	}
	        for m := 0; m < t2.Size.Y; m++ {     // kernel rows
	            mm := t2.Size.Y - 1 - m;      	// row index of flipped kernel
	            for n := 0; n < t2.Size.X; n++ {	// kernel columns

	                nn := t2.Size.X - 1 - n;  // column index of flipped kernel

	                // index of input signal, used for checking boundary
	                ii := i + (kCenterY - mm);
	                jj := j + (kCenterX - nn);

	                // ignore input samples which are out of bound
	                if ii >= 0 && ii < (*t1).Size.Y && jj >= 0 && jj < (*t1).Size.X {
	                	outputElement := (*t1).Data[ii + jj] * t2.Data[mm + nn]
	                    outputData.Data[el] += outputElement
	                    // fmt.Println(i, j, el, outputData.Data[el])
					}
	            }
	        }
	        fmt.Println(i, j, el, outputData.Data[el])
	    }
	}
	fmt.Print("Input Data: ", t1, "\n", "Kernel: ", t2, "\n", "Output Data: ", outputData)
	return outputData
}


func main() {
	inputData := NewTensor(7, 7)
	inputData.SetData(7, 7, []float64{1, 5, 6, 4, 5, 4, 4, 6, 1, 4, 4, 4, 2, 7, 2, 3, 9, 3, 4, 4, 4, 9, 8, 4, 6, 3, 2, 4, 6, 4, 3, 3, 2, 1, 3, 1, 1, 1, 1, 1, 1, 1, 2, 2, 3, 4, 2, 5, 4})
	//fmt.Println(inputData)
	//inputData.Print()

	kernel := NewTensor(3, 3)
	kernel.SetData(3, 3, []float64{1, 1, 1, 2, 2, 2, 3, 3, 3})
	//fmt.Println(kernel)
	//inputData.Print()
	stride := [2]int{2, 2}
	padding := [2]int{0, 0}
	inputData.Conv2D(kernel, stride, padding)
}