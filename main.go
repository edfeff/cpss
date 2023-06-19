package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_PASSWORD = "&#%345@789Abcdef"

type Options struct {
	Decrypt  bool
	Help     bool
	Password string
	Input    string
	Output   string
}

func main() {
	options := parseOption()
	if options.Help {
		showHelpInfo()
		return
	}
	checkOptions(options)

	process(options)

}

func checkOptions(options *Options) {
	if options.Password == "" {
		options.Password = DEFAULT_PASSWORD
	}
	if options.Input == "" {
		fmt.Println("请使用 -in 指定需要处理的文件")
		os.Exit(1)
	}
	if options.Output == "" {
		input := options.Input
		index := strings.LastIndex(input, ".")

		fileName := input
		suffix := ""

		if index != -1 {
			fileName = input[:index]
			suffix = input[index:]
		}
		now := time.Now().Nanosecond()
		if !options.Decrypt {
			options.Output = fileName + ".b" + strconv.Itoa(now) + suffix
		} else {
			options.Output = fileName + ".s" + strconv.Itoa(now) + suffix
		}
	}
}

func process(options *Options) {
	// abs path
	inFilePath := options.Input
	outFilePath := options.Output
	if !filepath.IsAbs(options.Input) {
		inFilePath, _ = filepath.Abs(options.Input)
	}
	if !filepath.IsAbs(outFilePath) {
		outFilePath, _ = filepath.Abs(options.Output)
	}

	// 设置AES密钥，长度必须是16、24或32字节
	password := padding(options.Password)
	key := []byte(password)

	if !options.Decrypt {
		// 加密文件
		err := encryptFile(key, inFilePath, outFilePath)
		if err != nil {
			fmt.Println("文件加密失败:", err)
			os.Exit(1)
		}
		fmt.Println("加密结果为", outFilePath)
	} else {
		// 解密文件
		err := decryptFile(key, inFilePath, outFilePath)
		if err != nil {
			fmt.Println("文件解密失败:", err)
			os.Exit(1)
		}
		fmt.Println("解密结果为", outFilePath)
	}

}

func padding(password string) string {
	if len(password) < 16 {
		n := 16 - len(password)
		for i := 0; i < n; i++ {
			password = password + "0"
		}
	}
	return password
}

func showHelpInfo() {
	help := `
	帮助信息:
	功能: 将文件加密成文本文件
	参数: 
			-d   添加此选项表示解密 
			-h   显示帮助信息
			-in  -i 文件
			-p   指定加密的密码
			-out -o 指定输出文件
	示例: 
	加密示例 cpss.exe -p 123456 -i xxx.mp4 -o xxx.mp4.bin
	解密示例 cpss.exe -p 123456 -i xxx.mp4.bin -o xxx.mp4

	`
	fmt.Println(help)
}

func parseOption() *Options {
	options := &Options{}
	flag.BoolVar(&options.Decrypt, "d", false, "encrypt")
	flag.StringVar(&options.Password, "p", "", "set password.")
	flag.StringVar(&options.Input, "in", "", "input file.")
	flag.StringVar(&options.Input, "i", "", "input file.")
	flag.StringVar(&options.Output, "o", "", "output file.")
	flag.StringVar(&options.Output, "out", "", "output file.")
	flag.BoolVar(&options.Help, "h", false, "show help message.")
	flag.Parse()
	return options
}

func encryptFile(key []byte, inputFile string, outputFile string) error {
	plainText, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// 创建一个用于加密的GCM模式的实例
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 生成随机的nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// 使用gcm.Seal加密数据
	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	// 在输出文件中写入nonce和密文
	err = ioutil.WriteFile(outputFile, append(nonce, cipherText...), 0644)
	if err != nil {
		return err
	}
	return nil
}

func decryptFile(key []byte, inputFile string, outputFile string) error {
	cipherText, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// 创建一个用于加密的GCM模式的实例
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 提取nonce和密文
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	// 使用gcm.Open解密数据
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return err
	}

	// 在输出文件中写入解密后的数据
	err = ioutil.WriteFile(outputFile, plainText, 0644)
	if err != nil {
		return err
	}
	return nil
}

// go run main.go -p 12345 -in main.go -out main.go.bin
// go run main.go -d -p 12345 -in main.go.bin -out main.go.txt
