package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	//苹果手机的路径
	src := "/Users/fuling/ali/sourcetree/go_copy_dir/testSrc"

	//windows 硬盘的路径
	dest := "/Users/fuling/ali/sourcetree/go_copy_dir/testTarget"

	//比较两个文件夹有什么不一样
	err := diffDir(src, dest)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//从苹果拷贝到 windows
	err = copyDir(src, dest)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func diffDir(src string, dest string) error {
	// 读取源目录
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	fmt.Println("原文件夹 mi：", src, ",共：", len(entries), "个文件")

	// 读取目标目录
	destEntries, err := ioutil.ReadDir(dest)
	if err != nil {
		return err
	}
	fmt.Println("目标文件夹 m2：", src, ",共：", len(destEntries), "个文件")

	data := make(map[string]int, 0)
	for _, entry := range entries {
		data[entry.Name()] = int(entry.Size())
	}

	desData := make(map[string]int, 0)
	for _, entry := range destEntries {
		desData[entry.Name()] = int(entry.Size())
	}
	printDiff(data, desData)
	return nil
}

func DifferenceOfKeys(m1, m2 map[string]int) (onlyInM1, onlyInM2, inBoth, inBothWithDifferenceValue map[string]int) {
	onlyInM1 = make(map[string]int)
	onlyInM2 = make(map[string]int)
	inBoth = make(map[string]int)
	inBothWithDifferenceValue = make(map[string]int)

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			inBoth[k1] = v1
			if v1 != v2 {
				inBothWithDifferenceValue[k1] = v1
			}
		} else {
			onlyInM1[k1] = v1
		}
	}

	for k2, v2 := range m2 {
		if _, ok := m1[k2]; !ok {
			onlyInM2[k2] = v2
		}
	}

	return
}

func printDiff(m1, m2 map[string]int) {

	onlyInM1, onlyInM2, _, inBothWithDifferenceValue := DifferenceOfKeys(m1, m2)

	fmt.Println("Only in m1:", onlyInM1)
	fmt.Println("Only in m2:", onlyInM2)
	//fmt.Println("In both:", inBoth)
	fmt.Println("In both with different values:", inBothWithDifferenceValue)
}
func copyDir(src string, dest string) error {
	// 读取源目录
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	fmt.Println("原文件夹：", src, ",共：", len(entries), "个文件")

	// 读取目标目录
	destEntries, err := ioutil.ReadDir(dest)
	if err != nil {
		return err
	}

	fmt.Println("目标文件夹：", dest, ",共：", len(destEntries), "个文件")

	if len(destEntries) > len(entries) {
		return errors.New("校验错误，目标文件夹的文件个数不应该 > 原文件夹的文件个数")
	}
	allSize := int64(0)
	for _, entry := range entries {
		//fmt.Println("正在拷贝第:", idx, "共", len(entries), "个文件需要拷贝")
		srcPath := filepath.Join(src, entry.Name())
		srcStat, err := os.Stat(srcPath)
		if err != nil {
			return err
		}
		allSize += srcStat.Size()
	}
	fmt.Println("原文件夹：", src, ",共：", allSize/1024/1024, "Mb")

	// 创建目标目录
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	success := 0
	start := time.Now()
	for idx, entry := range entries {
		//time.Sleep(time.Second)
		//fmt.Println("正在拷贝第:", idx, "共", len(entries), "个文件需要拷贝")
		srcPath := filepath.Join(src, entry.Name())
		srcStat, err := os.Stat(srcPath)
		if err != nil {

			return err // 其他类型的错误，例如权限问题
		}
		ms := time.Now().UnixMilli() - start.UnixMilli()

		successCount := (idx + 1) * 100 / len(entries)

		fmt.Println(time.Now(), " - 进度:", (idx+1)*100/len(entries), "% 已经消耗", ms/1000, "秒. 正在拷贝第:", idx+1, "个文件，共", len(entries), "个文件，文件名", entry.Name(), ",大小:", srcStat.Size()/1024/1024, "Mb")

		lessMs := int(ms) * (100 - successCount) / successCount

		fmt.Println("")
		fmt.Print("大约还需要", lessMs/1000, "秒，进度条[")
		for i := 0; i < successCount; i++ {
			fmt.Print(">")
		}
		for i := 0; i < 100-successCount; i++ {
			fmt.Print(".")
		}
		fmt.Println("]")
		fmt.Println("")
		destPath := filepath.Join(dest, entry.Name())

		// 判断是否是文件夹
		if entry.IsDir() {
			return errors.New("错误发送，不支持文件夹里还有文件夹" + srcPath)
			// 递归复制文件夹
			//err = copyDir(srcPath, destPath)
			//if err != nil {
			//	return err
			//}
		}

		destStat, err := os.Stat(destPath)
		if err != nil {
			//if os.IsNotExist(err) {
			//	//
			//	//return false, nil
			//}
			//return err // 其他类型的错误，例如权限问题
		} else {
			//存在
			fmt.Println("...***...13可能重名文件 Create fail，拷贝失败:", "destPath", destPath)

			if destStat.Size() != srcStat.Size() {
				fmt.Println("...***...3可能重名文件 文件大小还不一样，拷贝失败:", "destPath", destPath)
			}
			continue
		}

		output, err := os.Create(destPath)
		if err != nil {
			fmt.Println("...***...可能重名文件 Create fail，拷贝失败:", "destPath", destPath)
			continue
			//return err
		}
		defer output.Close()

		// 复制文件
		input, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer input.Close()
		_, err = io.Copy(output, input)
		if err != nil {
			fmt.Println("...***...运行错误 Copy fail，拷贝失败:", "destPath", destPath)
			continue
		}
		fmt.Println("...***... 拷贝成功:", "destPath", destPath)
		success++
	}

	fmt.Println("#####################  #########  #########  #########  #########  #########  #########    ")
	fmt.Println("############  拷贝报告:", "共成功拷贝", success, " 个文件,所有文件共", len(entries), " 个文件,拷贝率:", success*100/len(entries), "%")
	fmt.Println("#####################  #########  #########  #########  #########  #########  #########    ")

	return nil
}
