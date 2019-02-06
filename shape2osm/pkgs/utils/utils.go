package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  //
	_ "image/jpeg" //
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/lunixbochs/struc"
)

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func StructToBytes(class interface{}) []byte {
	var temp_buffer bytes.Buffer
	err := struc.Pack(&temp_buffer, class)
	CheckError(err)
	return temp_buffer.Bytes()
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func WriteError(fname string, text string, err error) {
	file, file_err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if file_err != nil {
		fmt.Println("Can't create or open log file:", fname)
	}
	defer file.Close()
	_, write_err := file.WriteString(fmt.Sprintf("Time: %v; Function: %v; Error: %v;\n", time.Now(), text, err))
	if write_err != nil {
		fmt.Println("Can't write into the log file:", fname)
	}
}

func QS_Int(my_slice *[]int, left, last int) {
	if left >= last {
		return
	}
	currleft := (*my_slice)[left]
	cab := left + 1
	for i := left; i <= last; i++ {
		if currleft > (*my_slice)[i] {
			(*my_slice)[cab], (*my_slice)[i] = (*my_slice)[i], (*my_slice)[cab]
			cab++
		}
	}
	(*my_slice)[left], (*my_slice)[cab-1] = (*my_slice)[cab-1], currleft
	QS_Int(my_slice, left, cab-2)
	QS_Int(my_slice, cab, last)
}

func QS_Float(my_slice *[]float64, left, last int) {
	if left >= last {
		return
	}
	currleft := (*my_slice)[left]
	cab := left + 1
	for i := left; i <= last; i++ {
		if currleft > (*my_slice)[i] {
			(*my_slice)[cab], (*my_slice)[i] = (*my_slice)[i], (*my_slice)[cab]
			cab++
		}
	}
	(*my_slice)[left], (*my_slice)[cab-1] = (*my_slice)[cab-1], currleft
	QS_Float(my_slice, left, cab-2)
	QS_Float(my_slice, cab, last)
}

func StringUnixToTime(unixString *string) time.Time {
	tmp_unix, err := strconv.ParseInt((*unixString), 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(tmp_unix, 0)
}

func StringInSlice(elem string, slice []string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

func StringInSlicePointer(elem *string, slice *[]string) bool {
	for _, e := range *slice {
		if e == (*elem) {
			return true
		}
	}
	return false
}

func Round(v float64) float64 {
	if v >= 0 {
		return math.Floor(v + 0.5)
	} else {
		return math.Ceil(v - 0.5)
	}
}

func RoundPlaces(v float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(v*shift) / shift
}

func MaxArrayFloat(v []float64) (float64, int) {
	var index int
	var max float64
	for i := range v {
		if v[i] > max {
			max = v[i]
			index = i
		}
	}
	return max, index
}

func MaxArrayInt(v []int) (int, int) {
	var index int
	var max int
	for i := range v {
		if v[i] > max {
			max = v[i]
			index = i
		}
	}
	return max, index
}

func ReverseByte(val byte) byte {
	var rval byte = 0
	for i := uint(0); i < 8; i++ {
		if val&(1<<i) != 0 {
			rval |= 0x80 >> i
		}
	}
	return rval
}

func ReverseUint8(val uint8) uint8 {
	return ReverseByte(val)
}

func ReverseUint16(val uint16) uint16 {
	var rval uint16 = 0
	for i := uint(0); i < 16; i++ {
		if val&(uint16(1)<<i) != 0 {
			rval |= uint16(0x8000) >> i
		}
	}
	return rval
}

//
func MakeInt64(v interface{}, ans *int64) error {
	var err error
	switch v.(type) {
	case int64:
		(*ans) = v.(int64)
		return err
	default:
		err = errors.New("Can't make int64 from interface")
		return err
	}
}

func MakeFloat64(v interface{}, ans *float64) error {
	var err error
	switch v.(type) {
	case float64:
		(*ans) = v.(float64)
		return err
	default:
		err = errors.New("Can't make float64 from interface")
		return err
	}
}

func MakeString(v interface{}, ans *string) error {
	var err error
	switch v.(type) {
	case string:
		(*ans) = v.(string)
		return err
	default:
		err = errors.New("Can't make string from interface")
		return err
	}
}

func ScaleBase64(base64string *string, width int, height int) error {
	var err error
	if (*base64string) == "" {
		err = errors.New("Empty base64 string")
		return err
	}
	fromBase64, err := base64.StdEncoding.DecodeString(*base64string)
	if err != nil {
		return err
	}
	r := bytes.NewReader(fromBase64)
	myImage, err := png.Decode(r)
	myImage = imaging.Resize(myImage, width, height, imaging.Box)
	var buff bytes.Buffer
	png.Encode(&buff, myImage)
	*base64string = base64.StdEncoding.EncodeToString(buff.Bytes())
	return err
}

func RotateBase64(base64string *string, angle float64) error {
	var err error
	if (*base64string) == "" {
		err = errors.New("Empty base64 string")
		return err
	}
	fromBase64, err := base64.StdEncoding.DecodeString(*base64string)
	if err != nil {
		return err
	}
	r := bytes.NewReader(fromBase64)
	myImage, err := png.Decode(r)
	img_rotate := imaging.Rotate(myImage, (angle * 180 / math.Pi), color.Transparent)
	var buff bytes.Buffer
	png.Encode(&buff, img_rotate)
	*base64string = base64.StdEncoding.EncodeToString(buff.Bytes())
	return err
}

func CustomRotate(img *image.Image, angle float64) *image.NRGBA {
	b := (*img).Bounds()
	imgOUT := image.NewNRGBA(image.Rect(0, 0, b.Max.X, b.Max.Y))
	//angle := math.Pi / 4
	for x := 0; x < b.Max.X; x++ {
		for y := 0; y < b.Max.Y; y++ {
			x0 := float64(b.Max.X) / 2
			y0 := float64(b.Max.Y) / 2
			X := int(Round(x0 + (float64(x)-x0)*math.Cos(angle) - (float64(y)-y0)*math.Sin(angle)))
			Y := int(Round(y0 + (float64(y)-y0)*math.Cos(angle) + (float64(x)-x0)*math.Sin(angle)))

			if X > 0 && X < b.Max.X && Y > 0 && Y < b.Max.Y {
				imgOUT.Set(X, Y, (*img).At(x, y))
			}
		}
	}
	return imgOUT
}

func ConvertToString(variable interface{}, tmp *string) error { //*
	var err error
	switch variable.(type) {
	case string:
		*tmp = variable.(string)
		return err

	default:
		err = errors.New("Error! String expexted")
		return err
	}
}

func ConvertToInt64(variable interface{}, tmp *int64) error {
	var err error
	switch variable.(type) {
	case int64:
		*tmp = variable.(int64)
		return err
	default:
		err = errors.New("Error! Int64 expected")
		return err
	}
}

func ConvertToInt(variable interface{}, tmp *int) error {
	var err error
	switch variable.(type) {
	case int:
		*tmp = variable.(int)
		return err
	default:
		err = errors.New("Error! Int expected")
		return err
	}
}

func ConvertToBool(variable interface{}, tmp *bool) error {
	var err error
	switch variable.(type) {
	case bool:
		*tmp = variable.(bool)
		return err
	default:
		err = errors.New("Error! Bool expected")
		return err
	}
}

func NowDateString() string {
	date := time.Now().Format("2006_01_02")
	return date
}

func NowTimeString() string {
	date := time.Now().Format("2006-01-02 15:04:05")
	return date
}

func GetQeryIN(district string) string {
	regions := strings.Split(district, "|")
	var regQuery string
	for i := 0; i < len(regions); i++ {
		regQuery += "'" + regions[i] + "'"
		if i+1 != len(regions) {
			regQuery += ","
		}
	}
	return regQuery
}

// CreateINQuery Create 'IN("1","2","3",...)'
func CreateINQuery(arr []string) string {
	var ans = "IN("
	if len(arr) == 0 {
		return ""
	}
	for i := range arr {
		if i == len(arr)-1 {
			ans += fmt.Sprintf("'%v')", arr[i])
		} else {
			ans += fmt.Sprintf("'%v',", arr[i])
		}
	}
	return ans
}

// CreateINQueryInt Create 'IN(1,2,3,...)'
func CreateINQueryInt(arr []int) string {
	var ans = "IN("
	if len(arr) == 0 {
		return ""
	}
	for i := range arr {
		if i == len(arr)-1 {
			ans += fmt.Sprintf("'%v')", arr[i])
		} else {
			ans += fmt.Sprintf("'%v',", arr[i])
		}
	}
	return ans
}

func GetQeryLike(search string) string {
	data := strings.Split(search, "|")
	var regQuery string
	for i := 0; i < len(data); i++ {

		if i+1 != len(data) {
			regQuery += fmt.Sprintf("lower('%%%v%%'),", data[i])
		} else {
			regQuery += fmt.Sprintf("lower('%%%v%%')", data[i])
		}
	}
	return regQuery
}
