const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
	    port: '8080', // 设置端口号
	    proxy: {
	        '/v1': {
	          target: 'http://127.0.0.1:8000',
	          ws: true,
	          changeOrigin: true,
	          pathRewrite: {
	            '^/v1': '/v1',
	          }
	        }
	    },
	}
})
