<h1 align="center">
	webscan
</h1>

<h4 align="center">web信息收集工具</h4>

<p align="center">
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/license-MIT-_red.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/niudaii/webscan">
    <img src="https://goreportcard.com/badge/github.com/niudaii/webscan?style=flat-square">		
  </a>
  <a href="https://github.com/niudaii/webscan/actions">
    <img src="https://img.shields.io/github/workflow/status/niudaii/webscan/Release?style=flat-square" alt="Github Actions">
  </a>
  <a href="https://github.com/niudaii/webscan/releases">
    <img src="https://img.shields.io/github/release/niudaii/webscan/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/niudaii/webscan/releases">
  	<img src="https://img.shields.io/github/downloads/niudaii/webscan/total">
  </a>
</p>


## 功能

- 
- 支持彩色输出
- 支持多种输出模式（debug|silent）
- 全平台支持
- API模式，可参考（internal/runner/runner.go）

## 使用

```
➜  webscan ./webscan -h
Webscanner

Usage:
  ./webscan [flags]

Flags:
INPUT:
   -i, -input string       url input(example: -i 'http://www.baidu.com', -i '192.168.243.11:81')
   -f, -input-file string  urls file(example: -f 'xxx.txt')

CONFIG:
   -threads int         number of threads (default 1)
   -timeout int         timeout in seconds (default 10)
   -p, -proxy string    proxy(example: -p 'http://127.0.0.1:8080')
   -header string       add custom headers(example: -header 'User-Agent: xxx, ')
   -finger-file string  use your finger file(example: -finger-file 'fingers.json')

OUTPUT:
   -o, -output string  output file to write found results (default "webscan.txt")
   -nc, -no-color      disable colors in output

DEBUG:
   -silent  show only results in output
   -debug   show debug output
```



## 参考

https://github.com/Becivells/iconhash