package handle

import (
	"bytes"
	"encoding/hex"
	"fileserver/public"
	logs "fileserver/public/plugin/logs"
	"fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

type Result struct {
	ResponseRes []public.UploadResponse
}

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif")  //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

func Ppload(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
	writer.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type,Auth,Authorization")
	writer.Header().Add("Access-Control-Allow-Credentials", "true")
	if request.Method != http.MethodPost {
		writer.Write(public.GetHttpStringMessage(public.CodeWrongHttpMethod_1001))
		return
	}
	if err := request.ParseMultipartForm(32 << 25); err != nil {
		logs.Fatal(err)
		writer.Write(public.GetHttpStringMessage(public.CodeInternalWrong_500))
		return
	}
	file := request.MultipartForm.File["file"]

	if file == nil {
		logs.Fatal("file is null")
		writer.Write(public.GetHttpStringMessage(public.CodeInternalWrong_500))
	}
	var cs Result
	for i := 0; i < len(file); i++ {
		upFile, err := file[i].Open()
		defer upFile.Close()

		path := util.GetUuid(file[i].Filename) + path.Ext(file[i].Filename)
		fSrc, err := ioutil.ReadAll(upFile)         //  使用ioutil.ReadAll(file)  会把文件读取掉  从而导致io.copy()  的结果是空文件
		res := GetFileType(fSrc[:10])
		if res==""{
			logs.Fatal(fmt.Sprintf("%s is not image",file[i].Filename))
			writer.Write(public.GetHttpStringMessage(public.CodeFileNotExist_2))
			return
		}

		cur, err := os.Create(public.Dir + string(os.PathSeparator) + path)
		defer cur.Close()
		if err != nil {
			logs.Fatal(err)
			writer.Write(public.GetHttpStringMessage(public.CodeInternalWrong_500))
			return
		}
		if _, err := io.Copy(cur,upFile); err != nil {
			logs.Fatal(err)
			writer.Write(public.GetHttpStringMessage(public.CodeInternalWrong_500))
			return
		}
		cs.ResponseRes = append(cs.ResponseRes,public.UploadResponse{public.CommonResponse{public.CodeSuccess_0, public.GetMessage(public.CodeSuccess_0)},
			public.UploadResponseBody{"/backend/file/pic/" + path}})

	}
	writer.Write(public.GetHttpByteBody(cs.ResponseRes))
	return
}
