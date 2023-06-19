# cpss

一个GO编写的文件加解密器

```text
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
```