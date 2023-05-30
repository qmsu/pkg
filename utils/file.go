package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// 判断文件夹是否存在  true--存在
func DirIsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 获取安装程序的根目录
func GetRootPath() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir + `\`
}

// 获取父文件夹
func getParentDirectory(directory string) string {
	return directory[0:strings.LastIndex(directory[0:len(directory)-1], "\\")]
}

const (
	IgnoreType_None       = iota //不过滤
	IgnoreType_Suffix            //后缀过滤
	IgnoreType_Prefix            //前缀过滤
	IgnoreType_FileName          //文件名称过滤
	IgnoreType_SelfIgnore        //自定义过滤
)

type Ignore struct {
	IgnoreContent string
	Type          int
	IgnoreFnc     func(ignoreContent string) bool
}

//删除空文件夹
func DelDir(rootPath, filePath string) {
	dirPath := rootPath + filePath
	for j := 0; j < len(strings.Split(rootPath+filePath, "\\")); j++ {
		nullIgnoreMap := make(map[int]Ignore)
		if DirIsExist(dirPath) && CountFile(dirPath, nullIgnoreMap) == 0 {
			_ = os.Remove(dirPath)
		}
		dirPath = getParentDirectory(dirPath)
		if dirPath+"\\" == rootPath {
			break
		}
	}
}
func IsIgnoreFile(fileName string, ignoreMap map[int]Ignore) bool {
	isIgnore := false
	for ignoreType, ignore := range ignoreMap {
		if ignoreType == IgnoreType_Suffix {
			if strings.HasSuffix(fileName, ignore.IgnoreContent) {
				isIgnore = true
			}
		} else if ignoreType == IgnoreType_Prefix {
			if strings.HasPrefix(fileName, ignore.IgnoreContent) {
				isIgnore = true
			}
		} else if ignoreType == IgnoreType_FileName {
			if strings.Contains(fileName, ignore.IgnoreContent) {
				isIgnore = true
			}
		} else if ignoreType == IgnoreType_SelfIgnore && ignore.IgnoreFnc != nil {
			if ignore.IgnoreFnc(ignore.IgnoreContent) {
				isIgnore = true
			}
		}
	}
	return isIgnore
}

func CountFile(dir string, ignoreMap map[int]Ignore) int {
	var count int
	_ = filepath.Walk(dir, func(oldPath string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if IsIgnoreFile(oldPath, ignoreMap) {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".temp") || strings.HasSuffix(info.Name(), ".cp") {
			return nil
		}
		fileArr := strings.Split(info.Name(), ".")
		if len(fileArr) >= 2 {
			if fileArr[0] != "" {
				count++
			}
		}
		return nil
	})
	return count
}

// 判断文件目录是否存在，不存在即创建该文件夹
func CreatePathIfNotExists(_dir string) (bool, error) {
	_, dirExist := os.Stat(_dir)
	if dirExist == nil {
		return true, nil
	} else {
		// 创建文件夹
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			return true, err
		} else {
			return true, nil
		}
	}
}

//列出指定目录下的目录
func ListDirDir(dir string) ([]string, error) {
	s, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}
	var res []string
	for _, v := range s {
		if IsFile(v) {
			continue
		}
		res = append(res, v)
	}
	return res, nil
}

//拷贝文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	//获取源文件的权限
	fi, _ := src.Stat()
	perm := fi.Mode()

	dst, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm) //复制源文件的所有权限
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//重命名文件
//如果重复则重命名为 文件名 - 副本 的形式
func RenameFile(src string) (newName string, err error) {
	info, err := os.Stat(src)
	if err == nil {
		i := 1
		for {
			var tmp string
			if info.IsDir() {
				newName = fmt.Sprintf("%s - 副本", newName)
				tmp = src + newName
			} else {
				tmp = filepath.Dir(src)
				ext := filepath.Ext(info.Name())
				arr := strings.Split(info.Name(), ext)
				if newName == "" {
					newName = arr[0]
				}
				newName = fmt.Sprintf("%s - 副本", newName)
				tmp = filepath.Join(tmp, newName+ext)
			}
			_, err = os.Stat(tmp)
			if err == nil {
				i++
				continue
			}
			if os.IsNotExist(err) {
				return tmp, nil
			} else {
				return "", err
			}
		}
	}
	if os.IsNotExist(err) {
		return src, nil
	}
	return "", err
}

//重命名文件
//如果重复则重命名为 文件名_序号 的形式
func RenameFileName(targetPath, fileName, fileExt string) (newName string, err error) {
	newPath := filepath.Join(targetPath, fileName+fileExt)
	_, err = os.Stat(newPath)
	if err == nil {
		i := 1
		for {
			//判断文件是否存在，如果存在则重命名为文件名_序号
			newName = fmt.Sprintf("%s_%02d%s", fileName, i, fileExt)
			_, err = os.Stat(filepath.Join(targetPath, newName))
			if err == nil {
				i++
				continue
			}
			if os.IsNotExist(err) {
				return newName, nil
			} else {
				return "", err
			}
		}
	}
	if os.IsNotExist(err) {
		return fileName + fileExt, nil
	}
	return "", err
}

func GetFileExt(file string) string {
	ext := filepath.Ext(filepath.Base(file))
	arr := strings.Split(ext, ".")
	if len(arr) == 2 {
		return strings.ToUpper(arr[1])
	}
	return ""
}

//获取文件后缀，返回结果中不带.
func GetFileExtWithoutPoint(file string) string {
	fileName := filepath.Base(file)
	index := strings.LastIndex(fileName, ".")
	if index == -1 {
		return ""
	}
	ext := fileName[index+1:]
	return strings.ToLower(ext)
}

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("474946383961", "gif")          //GIF (gif)
	fileTypeMap.Store("474946383761", "gif")          //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)
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

//判断照片后缀和格式是否匹配
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

func GetRealFile(pic string) (string, error) {
	var newXJPicPath string
	ext, err := GetImgExt(pic)
	if err != nil {
		return pic, err
	}
	if ext != "" && "."+ext != filepath.Ext(pic) {
		arr := strings.Split(filepath.Base(pic), ".")
		if len(arr) != 2 {
			return pic, err
		}
		newXJPicPath = filepath.Join(filepath.Dir(pic), arr[0]+"."+ext)
		err = os.Rename(pic, newXJPicPath)
		if err != nil {
			return pic, err
		}
	} else {
		newXJPicPath = pic
	}
	return newXJPicPath, nil
}

func GetImgExt(pic string) (string, error) {
	file, err := os.Open(pic)

	if err != nil {
		return "", err
	}

	defer file.Close()

	buff := make([]byte, 512)

	_, err = file.Read(buff)

	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	if strings.HasPrefix(filetype, "image/") {
		return strings.TrimPrefix(filetype, "image/"), nil
	}

	return "", err
}

//获取某个目录所在的磁盘类型
func GetDirFType(dir string) (string, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return "", err
	}
	for _, part := range parts {
		if strings.HasPrefix(dir, part.Mountpoint) {
			return part.Fstype, nil
		}
	}
	return "", err
}

// 去除文件路径中的制表符回车等影响文件创建的因素
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("[\f\n\r\t\v]")
	return reg.ReplaceAllString(str, "")
}
