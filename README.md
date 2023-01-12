## CexBot 介绍

本程序提供 Binance CoinBase Gate.io Huobi KuCoin MEXC OKX的上币监控

telegram频道:https://t.me/zcchann_Blockchain

## 设置

config.json中需要填入以下内容:

```json
{
  "binance_config": {
    "api_key": "",   
    "secret_key": "" 
  },
  "okx_config": {
    "api_key": "",
    "secret_key": "",
    "passphrase": ""
  },
  "mongodb": {
    "username": "",
    "password": "",
    "host": "",
    "port": "",
    "database": ""
  },
  "telegram": {
    "bot_id": "",            #telegram 机器人ID
    "channel_id": "",        #推送上币信息的频道ID
    "error_channel_id": ""   #推送报错信息的频道ID
  }
}
```

### 注意

推送时如果日志返回错误信息：

{"ok":false,"error_code":400,"description":"Bad Request: message text is empty"}

需要在channel_id前加100

例：

telegram的channel_id为 "-123456789"

在config.json中需要填写为："channel_id": "-100123456789"

## 关于数据库

程序对mongoDB数据库版本没有要求

在shell文件夹中我们提供了一个简易的mongo安装脚本

执行以下命令即可安装

```
chmod +x shell/install.sh
./shell/install.sh
```



如果数据库中只存放cexbot的数据 并且你对数据库中的数据权限比较无所谓 在安装完mongoDB后执行以下命令创建admin用户

```
[root@localhost ~]# mongo

> use admin

> db.createUser({ user:'admin', pwd:'123456',roles:[{role:'root',db:'admin'}]})
```



## 捐赠

如果你觉得这个程序对你有帮助,你可以通过下面的地址捐赠更好的支持cexbot

ETH address: **0x6aB9886a6B86F649F9FfAA97074F7aa2F461494E**