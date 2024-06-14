import tornado.ioloop
import tornado.httpserver 
import tornado.web 
import tornado.options 
from tornado.options import define, options
import requests
import logging
import json
logging.basicConfig(level = logging.INFO,format = '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
log = logging.getLogger(__name__)
define("port", type=int, default=8001, help="run on the given port")
# 创建请求处理器
# 当处理请求时会进行实例化并调用HTTP请求对应的方法


class IndexHandler(tornado.web.RequestHandler):
    def get(self):
        url = "http://demo-edge:8000/"
        res = requests.get(url=url)
        res = json.loads(res.text)
        number = res.get('number')
        logging.info(number)
        self.write({'number':number})
        self.finish()

# 创建路由表
urls = [(r"/", IndexHandler)
        ]

# 定义服务器
def main():
    # 解析命令行参数
    tornado.options.parse_command_line()
    # 创建应用实例
    app = tornado.web.Application(urls)
    # 监听端口
    app.listen(options.port)
    # 创建IOLoop实例并启动
    tornado.ioloop.IOLoop.current().start()

# 应用运行入口，解析命令行参数
if __name__ == "__main__":
    # 启动服务器
    main()
