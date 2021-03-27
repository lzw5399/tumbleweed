/**
 * @Author: lzw5399
 * @Date: 2021/3/27 16:47
 * @Desc: 应用启动时初始化相关的依赖
 */
package initialize

func init() {
	// logger
	setupLogger()

	// 配置
	setupConfig()

	// 数据库连接
	setupDbConn()

	// 内存缓存
	setupCache()
}
