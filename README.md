
##### 表格
|xx|xx|
|---|---|
| xx | xx |

###### dig容器
```
go run ./digtest/main.go -c ./digtest/my.ini

================== redis section =====================
redis ip: 127.0.0.1
redis port: 6379
redis db: 0
================== mysql section =====================
mysql ip: 127.0.0.1
mysql port: 3306
mysql user: root
mysql password: 123456
mysql db: test
~~~~~~~~~~~~~~~~~~ redis section connect ~~~~~~~~~~~~~
redis ip: 127.0.0.1
redis port: 6379
redis db: 0
```

###### Step1：编写描述文件：hello.proto
```
syntax = "proto3"; // 指定proto版本
package hello;     // 指定默认包名
// 指定golang包名
option go_package = "./;hello";
// 定义Hello服务
service Hello {
    // 定义SayHello方法
    rpc SayHello(HelloRequest) returns (HelloResponse) {}
}
// HelloRequest 请求结构
message HelloRequest {
    string name = 1;
}
// HelloResponse 响应结构
message HelloResponse {
    string message = 1;
}
```


###### Step2：编译生成.pb.go文件
```
$ cd proto/hello
# 编译hello.proto
$ protoc -I . --go_out=plugins=grpc:. ./hello.proto

```


###### Step3：实现服务端接口 server/main.go
```
服务端引入编译后的proto包，定义一个空结构用于实现约定的接口，接口描述可以查看hello.pb.go文件中的
HelloServer接口描述。实例化grpc Server并注册HelloService，开始提供服务。

运行：
$  go run .\hello\server\main.go
Listen on 127.0.0.1:50052  //服务端已开启并监听50052端口
```


###### Step4：实现客户端调用 client/main.go
```
客户端初始化连接后直接调用hello.pb.go中实现的SayHello方法，即可向服务端发起请求，
使用姿势就像调用本地方法一样。

运行：
$  go run .\hello\client\main.go
Hello gRPC.    // 接收到服务端响应
```



###### Step5：证书制作 制作私钥 (.key)  
```
注：不能联通

@REM # Key considerations for algorithm "RSA" ≥ 2048-bit
@REM $ openssl genrsa -out server1.key 2048
@REM # Key considerations for algorithm "ECDSA" ≥ secp384r1
@REM # List ECDSA the supported curves (openssl ecparam -list_curves)
@REM $ openssl ecparam -genkey -name secp384r1 -out server1.key
@REM ----
@REM ###### 自签名公钥(x509) (PEM-encodings .pem|.crt)
@REM $ openssl req -new -x509 -sha256 -key server1.key -out server1.pem -days 3650

自定义信息
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:XxXx
Locality Name (eg, city) []:XxXx
Organization Name (eg, company) [Internet Widgits Pty Ltd]:XX Co. Ltd
Organizational Unit Name (eg, section) []:Dev
Common Name (e.g. server FQDN or YOUR name) []:wuzhi555.cc
Email Address []:xxx@xxx.com

```


###### go-grpc TSL认证 解决
```
一、问题描述：transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable
二、背景环境：我的环境windows go 1.17，linux解决这个问题办法同样也适用。
```
三、首先需要 [下载SSL](http://slproweb.com/products/Win32OpenSSL.html)  
![xx](https://github.com/xiaowuzhi/hellov1/d1.png)

```

你点开这个链接以后会看到上面这图片显示的页面，我第一次下载的时候看见有个博客说随便点击一个下载，然后我也没仔细看，因为我的电脑是64位的，我就随便点了一个Win64的，然后好家伙，后面一共卸载下载了三次，不要下载forRAM的，会不允许安装，（ARM64是ARM中64位体系结构，x64是x86系列中的64位体系。ARM属于精简指令集体系，汇编指令比较简单。x86属于复杂指令集体系，汇编指令较多。属于两种不同的体系。看不懂没关系，你只要知道是两种不同的体系，那当然下载了也用不了）不要下载Light的，因为你会找不到后面所需要的openssl.cnf文件。
1、直接根据你的系统去下载最大的那两个其中之一。下载完成以后直接点开exe一直next安装好就可以了。
2、将openSSL的bin目录所在的路径放到path环境变量中，然后重启电脑。
3、生成普通的key
@REM openssl genrsa -des3 -out server.key 2048
@REM （记住设置的密码，命令直接在终端上执行就好，我直接在goland的终端上执行的）
无密码
openssl genrsa -out server.key 2048

4、生成ca的crt
openssl req -new -x509 -key server.key -out ca.crt -days 3650
遇到填东西的直接回车就行
5、生成csr
openssl req -new -key server.key -out server.csr
6、更改openssl.cnf （Linux 是openssl.cfg）
1）复制一份你安装的openssl的bin目录里面的openssl.cnf 文件到你项目所在的目录，我放在了keys文件夹下。
2）找到 [ CA_default ]，打开 copy_extensions = copy （就是把前面的#去掉）
3）找到[ req ]，打开 req_extensions = v3_req # The extensions to add to a certificate request
4）找到[ v3_req ]，添加 subjectAltName = @alt_names
5）添加新的标签 [ alt_names ]，和标签字段

DNS.1 = *.grpc.wuzhi555.cc
DNS.2 = *.wuzhi555.cc

7、生成证书私钥test.key
openssl genpkey -algorithm RSA -out test.key
8、通过私钥test.key生成证书请求文件test.csr（注意cfg和cnf）

openssl req -new -nodes -key test.key -out test.csr -days 3650  -config ./openssl.cnf -extensions v3_req

test.csr是上面生成的证书请求文件。ca.crt/server.key是CA证书文件和key，用来对test.csr进行签名认证。这两个文件在第一部分生成。
9、生成SAN证书

openssl x509 -req -days 365 -in test.csr -out test.pem -CA ca.crt -CAkey server.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req


10、然后就可以用在 GO 1.15 以上版本的GRPC通信了
服务器加载代码

creds, err := credentials.NewServerTLSFromFile("test.pem", "test.key")

客户端加载代码

creds,err := credentials.NewClientTLSFromFile("test.pem","*.wuzhi555.cc")


这个问题怎么说呢，试了很多的方法，最终使用这个方法解决了在这里记录一下。

```
















