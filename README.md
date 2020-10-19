# tencentcloud_infra_map_neo4j

## Usage

### What's this

- A simple demo to express resources on tencentcloud using neo4j graph datebase, shown as below:

![resource-region-map](pics/graph.png)

### How to run the project

1. Install neo4j server, could follow [neo4j-install-guide](https://blog.csdn.net/huacha__/article/details/81123410) 

2. Build a binary file from this project then run it.

### Note

1. Ak and sk needed to be exported to env before you run the binary file, using command below:

```bash
export TENCENTCLOUD_SECRET_ID=example
export TENCENTCLOUD_SECRET_KEY=example
```

2. If face an issue like `can't start neo4j service because of the tls`, try below steps:

```
1. 生成私钥
`openssl genrsa -out private.key 1024`
2. 生成证书
`openssl req -new -x509 -key private.key -out public.crt`
3. [配置TSL](https://s0neo4j0com.icopy.site/docs/operations-manual/3.4/security/ssl-framework/#term-ssl-certificate)
4. 生成私钥时不能带密码
```

3. If face an issue on `bolt connection` side, change `dbms.connector.bolt.tls_level` value to `OPTIONAL` in neo4j config file.

## Supporting products

### Products supported list:

- Region
- CVM
- Image

### For new products 

For new products wanted to be supported, could add it easily (I think), just put them into `product` directory.

## Acknowledgement

Thanks to [aws_infra_map_neo4j](https://github.com/rdkls/aws_infra_map_neo4j) for inspiring this idea.

## To be continued for more optimizing.