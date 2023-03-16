package utils

import (
	"MyLNPU/internal/logger"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
)

func GetVerifyCode(url string) (string, *http.Client, error) {
	client, err := downloadVerifyCode(url)
	if err != nil {
		return "", nil, err
	}
	codeImg, err := os.Open("./code.jpg")
	if err != nil {
		logger.Println("打开验证码失败...")
		return "", nil, err
	}
	img, _, err := image.Decode(codeImg)
	if err != nil {
		return "", nil, err
	}
	if img.Bounds().Empty() {
		return "", nil, image.ErrFormat
	}
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	var imgArr [20][72]int
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			colorRgb := img.At(j, i)
			_, g, _, _ := colorRgb.RGBA()
			gUint8 := uint8(g >> 8)
			if gUint8 < 150 {
				gUint8 = 0
				imgArr[i][j] = 0
			} else {
				gUint8 = 255
				imgArr[i][j] = 1
			}
		}
	}
	removeByLine(&imgArr) //去除杂线
	//toPrint(imgArr) //打印图片到控制台
	arr := cutImgArr(imgArr) //切割图像数组
	code := matchCode(arr)
	return code, client, nil
}

/*
匹配识别验证码
*/
func matchCode(imgArrArr [4][20][18]int) string {
	code := ""
	charMap := GetCharMap()
	for i := 0; i < len(imgArrArr); i++ {
		temp := ""
		maxMach := 0
		for k, v := range charMap {
			percent := compareStr(v, toString(imgArrArr[i]))
			if percent > maxMach {
				maxMach = percent
				temp = k
			}
		}
		code += temp
	}
	return code
}

/*
下载验证码图片
*/
func downloadVerifyCode(uri string) (*http.Client, error) {
	client, err := NewHttpClient()
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest("GET", uri, nil)
	get, err := client.Do(request)
	if err != nil {
		logger.Println("验证码获取失败...", err)
	}
	defer get.Body.Close()
	img, _, err := image.Decode(get.Body)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	rgba := image.NewRGBA(bounds)
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			colorRgb := img.At(j, i)
			_, g, _, a := colorRgb.RGBA()
			gUint8 := uint8(g >> 8)
			aUint8 := uint8(a >> 8)
			if gUint8 < 150 {
				gUint8 = 0
			} else {
				gUint8 = 255
			}
			rgba.SetRGBA(j, i, color.RGBA{R: gUint8, G: gUint8, B: gUint8, A: aUint8})
		}
	}
	subImage := rgba.SubImage(image.Rect(2, 11, 74, 31))
	create, _ := os.Create("./code.jpg")
	writer := bufio.NewWriter(create)
	err = jpeg.Encode(writer, subImage, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

/*
去除杂线
*/
func removeByLine(imgArr *[20][72]int) {
	for i := 1; i < 20-1; i++ {
		for j := 1; j < 72-1; j++ {
			if imgArr[i][j] == 0 {
				count := imgArr[i][j-1] + imgArr[i][j+1] + imgArr[i-1][j] + imgArr[i+1][j]
				if count > 2 {
					imgArr[i][j] = 1
				}
			}
		}
	}
}

/*
控制台打印
*/
func toPrint(imgArr [20][72]int) {
	for i := 0; i < len(imgArr); i++ {
		for j := 0; j < len(imgArr[i]); j++ {
			fmt.Printf("%d", imgArr[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

/*
控制台打印切割后的数组
*/
func cutArrToPrint(imgArrArr [4][20][18]int) {
	for i := 0; i < len(imgArrArr); i++ {
		for j := 0; j < len(imgArrArr[i]); j++ {
			for k := 0; k < len(imgArrArr[i][j]); k++ {
				fmt.Printf("%d", imgArrArr[i][j][k])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

/*
二维数组转字符串
*/
func toString(imgArr [20][18]int) string {
	var str string
	for i := 0; i < len(imgArr); i++ {
		for j := 0; j < len(imgArr[i]); j++ {
			str += strconv.Itoa(imgArr[i][j])
		}
	}
	return str
}

/*
对比两个字符串，返回相似度
*/
func compareStr(s1, s2 string) int {
	var percent int
	for i := 0; i < len(s1); i++ {
		if s1[i] == s2[i] {
			percent++
		}
	}
	return percent
}

/*
裁切
*/
func cutImgArr(imgArr [20][72]int) [4][20][18]int {
	var imgArrArr [4][20][18]int
	for i := 0; i < 20; i++ {
		for j := 0; j < 72; j++ {
			if j < 18 {
				t := j - 0*18
				imgArrArr[0][i][t] = imgArr[i][j]
			} else if j < 36 {
				t := j - 1*18
				imgArrArr[1][i][t] = imgArr[i][j]
			} else if j < 54 {
				t := j - 2*18
				imgArrArr[2][i][t] = imgArr[i][j]
			} else {
				t := j - 3*18
				imgArrArr[3][i][t] = imgArr[i][j]
			}
		}
	}
	return imgArrArr
}
