package config

type Config struct {
	PriKey string //私钥列表
	//ethrpc请求配置
	EthRpcConf struct {
		Url          string // Url:请求rpc地址 修改来替换各链rpc地址
		IntervalTime int    //请求轮询时间

	}
	//minted限制,用来过滤非热门铭文
	MintConf struct {
		ToAddr    string //发送地址
		InputData string //发送数据
	}
}
