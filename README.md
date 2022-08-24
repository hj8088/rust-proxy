# rust-proxy
类似nexus3的产品并没有很好的支持rust和cargo的代理支持，所以基于golang实现了一个简单的rust和cargo的代理，以满足在某一些内网环境下的代理功能
- 当前版本需要手动将rustup-init和crates.io-index的信息从外网导到内网环境下
- 通过类似的nexus3的工具，将crates和rust-static的远程服务通过raw的方式配置到内网

## Install

```
go get github.com/hj8088/rust-proxy
```

## Usage
```

- 手动下载crates的索引信息 crates.io-index ,并放置到代理程序的工作根目录下，索引目录的名称需要固定为'crates.io-index':   
    git clone git://mirrors.ustc.edu.cn/crates.io-index
    
- 更新creates的代理地址在 'config.json' 文件中:
    {
    "dl": "http://localhost:8080/api/v1/crates",
    "api": "https://crates.io/"
    }
- 并提交到本地的git中以生效:
    git add config.json
    git commit -m "update config.json"
- 通过nexus3将crates的目录和rust-static服务通过raw代理到内网     
    
- 启动:
    rust-proxy --server-port=8080 --project-root=/data/cargo-mirror --remote-crates-url=https://crates-io.proxy.ustclug.org/api/v1/crates --remote-rust-static-url=https://mirrors.ustc.edu.cn/rust-static
    
- 更新rust和cargo的代理：
    export RUSTUP_DIST_SERVER=http://localhost:8080/rust-static
    export RUSTUP_UPDATE_ROOT=http://localhost:8080/rust-static/rustup
 
  ~/.cargo/config
  [source.local]
  registry = "http://localhost:8080/crates.io-index"

```
