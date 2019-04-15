package utils

//Version 3.1

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const datefm = "20060102" //20060102150405 固定时间格式 Jan 2 15:04:05 2006 MST : 1 2 3 4 5 6 -7

func CheckPath(dir string) {
	if f, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	} else {
		if !f.IsDir() {
			panic(dir + " is file")
		}
	}
}

func GetColumnS(lst [][]string, c int) (vlst []string) {
	for _, v := range lst {
		vlst = append(vlst, v[c])
	}
	return
}

func GetColumnI(lst [][]string, c int) (vlst []int) {
	for _, v := range lst {
		if f, err := strconv.Atoi(v[c]); err == nil {
			vlst = append(vlst, f)
		} else {
			panic(err)
		}
	}
	return
}

func GetColumnF(lst [][]string, c int) (vlst []float64) {
	for _, v := range lst {
		if f, err := strconv.ParseFloat(v[c], 64); err == nil {
			vlst = append(vlst, f)
		} else {
			panic(err)
		}
	}
	return
}

func GetEMA(lst []float64, signal float64, TOLIST bool) (emalst []float64, ema float64) {
	k := 2.0 / (signal + 1.0)
	ema = lst[0]
	emalst = make([]float64, 0)

	if TOLIST {
		for _, x := range lst {
			ema += k * (x - ema)
			emalst = append(emalst, ema)
		}
	} else {
		for _, x := range lst {
			ema += k * (x - ema)
		}
	}
	return
}

func GetDifMacd(lst []float64) (diflst, macdlst []float64) {
	k9 := 2.0 / (9.0 + 1.0)
	k12 := 2.0 / (12.0 + 1.0)
	k26 := 2.0 / (26.0 + 1.0)
	ema12, ema26 := lst[0], lst[0]
	dea9 := 0.0

	for _, x := range lst {
		ema12 += k12 * (x - ema12)
		ema26 += k26 * (x - ema26)
		diflst = append(diflst, ema12-ema26)
		dea9 += k9 * (ema12 - ema26 - dea9)
		macdlst = append(macdlst, ema12-ema26-dea9)
	}
	return

}

func MinusList(lst1, lst2 []float64) (lst []float64) {
	for i, x := range lst1 {
		lst = append(lst, x-lst2[i])
	}
	return
}

func DifMacd(dif, macd []float64) (dmlst []int) {
	var p, n int
	for i := 0; i < len(dif); i++ {
		if dif[i] > 0 && macd[i] > 0 {
			p = 1
		} else {
			p = 0
		}
		if dif[i] < 0 && macd[i] < 0 {
			n = -1
		} else {
			n = 0
		}
		dmlst = append(dmlst, p+n)
	}
	return
}

func Fzero(a, b float64) float64 {
	if b <= 0 || a <= 0 {
		return -1
	} else {
		return math.RoundToEven(100*a/b) / 100
	}
}

func Addlst(lst []int, index []int) (sum int) {
	if index != nil {
		for _, i := range index {
			sum += lst[i]
		}
	} else {
		for _, x := range lst {
			sum += x
		}
	}
	return
}

func SumlstI(lst []int, index []int) (sum int) {
	if index != nil {
		for _, i := range index {
			sum += lst[i]
		}
	} else {
		for _, x := range lst {
			sum += x
		}
	}
	return
}

func SumlstF(lst []float64, index []int) (sum float64) {
	if index != nil {
		for _, i := range index {
			sum += lst[i]
		}
	} else {
		for _, x := range lst {
			sum += x
		}
	}
	return math.RoundToEven(100*sum) / 100
}

func MaxlstF(lst []float64, index []int) float64 {
	if index != nil {
		var l1 []float64
		for _, i := range index {
			l1 = append(l1, lst[i])
		}
		sort.Float64s(l1)
		return l1[len(l1)-1]
	} else {
		sort.Float64s(lst)
		return lst[len(lst)-1]
	}
}

func MinlstF(lst []float64, index []int) float64 {
	if index != nil {
		var l1 []float64
		for _, i := range index {
			l1 = append(l1, lst[i])
		}
		sort.Float64s(l1)
		return l1[0]
	} else {
		sort.Float64s(lst)
		return lst[0]
	}
}

func Translate(lst *[][]string) (tlst [][]string) {
	var maxline int
	for _, v := range *lst {
		if len(v) > maxline {
			maxline = len(v)
		}
	}

	var datalst [][]string
	for _, v := range *lst {
		if len(v) == maxline {
			datalst = append(datalst, v)
		} else {
			vv := append(v, make([]string, maxline-len(v))...)
			datalst = append(datalst, vv)
		}
	}

	var xlst []string
	for x := 0; x < maxline; x++ {
		xlst = nil
		for y := 0; y < len(datalst); y++ {
			xlst = append(xlst, datalst[y][x])
		}
		tlst = append(tlst, xlst)
	}

	return
}

func ArraySort(lst [][]string, firstIndex int, descending bool) [][]string {
	if len(lst) <= 1 {
		return lst
	}

	if firstIndex < 0 || firstIndex > len(lst[0])-1 {
		fmt.Println("Warning: Param firstIndex should between 0 and len(lst)-1. The original array is returned.")
		return lst
	}

	mLstArray := &LstArray{lst, firstIndex}
	if descending {
		sort.Sort(sort.Reverse(mLstArray))
	} else {
		sort.Sort(mLstArray)
	}
	return mLstArray.mArr
}

type LstArray struct {
	mArr       [][]string
	firstIndex int
}

func (arr *LstArray) Len() int {
	return len(arr.mArr)
}

func (arr *LstArray) Swap(i, j int) {
	arr.mArr[i], arr.mArr[j] = arr.mArr[j], arr.mArr[i]
}

func (arr *LstArray) Less(i, j int) bool {
	arr1 := arr.mArr[i]
	arr2 := arr.mArr[j]

	for index := arr.firstIndex; index < len(arr1); index++ {
		f1, err1 := strconv.ParseFloat(arr1[index], 64)
		f2, err2 := strconv.ParseFloat(arr2[index], 64)

		if err1 == nil && err2 == nil {
			v1 := f1
			v2 := f2
			if v1 < v2 {
				return true
			} else if v1 > v2 {
				return false
			}
		} else {
			v1 := arr1[index]
			v2 := arr2[index]
			if v1 < v2 {
				return true
			} else if v1 > v2 {
				return false
			}
		}

		//		if arr1[index] < arr2[index] {
		//			return true
		//		} else if arr1[index] > arr2[index] {
		//			return false
		//		}
	}
	return i < j
}

func Readline(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

func ReadFile(file string) [][]string {
	var datalst [][]string
	if by, err := ioutil.ReadFile(file); err == nil {
		for _, str := range strings.Split(strings.TrimRight(string(by), "\n"), "\n") {
			datalst = append(datalst, strings.Split(strings.TrimRight(str, "\r"), "\t"))
		}
	} else {
		//		panic(err)
	}

	first := 0
	last := len(datalst) - 1
	if len(datalst) > 0 {
		if _, err := strconv.Atoi(datalst[0][0]); err != nil {
			first += 1
		}

		if _, err := strconv.Atoi(datalst[last][0]); err != nil {
			last -= 1
		}

		if last < 0 || first >= len(datalst) {
			return nil
		}

		if len(datalst[first]) == 7 {
			if datalst[last][5] == "0" {
				if len(datalst) > 2 {
					last -= 1
				} else {
					return nil
				}
			}
		}
		return datalst[first : last+1]
	} else {
		return nil
	}
}

func Lst2Str(datalst []string) (datastr string) {
	for _, str := range datalst {
		datastr += str + "\t"
	}
	datastr = strings.TrimRight(datastr, "\t")
	return
}

func Lst2Str2(datalst [][]string) (datastr string) {
	for _, line := range datalst {
		for _, str := range line {
			datastr += str + "\t"
		}
		datastr = strings.TrimRight(datastr, "\t") + "\n"
	}
	datastr = strings.TrimRight(datastr, "\n")
	return
}

func ExportFileA(seekSize int64, data [][]string, expath string) {
	datastr := Lst2Str2(data)
	file, err := os.OpenFile(expath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	position, err2 := file.Seek(seekSize, 2)
	if err != nil {
		panic(err2)
	}
	file.Truncate(position)

	bufferedWriter := bufio.NewWriter(file)
	if _, err = bufferedWriter.WriteString(datastr); err != nil {
		panic(err)
	}
	bufferedWriter.Flush()
}

func ExportFile(dir, name string, data [][]string) {
	CheckPath(dir)
	date := time.Now().Format(datefm)
	expath := path.Join(dir, date+"-"+name)

	datastr := Lst2Str2(data)
	file, err := os.OpenFile(expath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, err = bufferedWriter.WriteString(datastr); err != nil {
		panic(err)
	} else {
		fmt.Println(expath)
	}
	bufferedWriter.Flush()
}

func ExportFileS(dir, name string, datastr string) {
	CheckPath(dir)
	date := time.Now().Format(datefm)
	expath := path.Join(dir, date+"-"+name)

	if del := os.Remove(expath); del != nil {
		fmt.Println("ExportFileS:", del)
	}

	file, err := os.OpenFile(expath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, err = bufferedWriter.WriteString(datastr); err != nil {
		panic(err)
	} else {
		fmt.Println(expath)
	}
	bufferedWriter.Flush()
}
