/*
go语言测试elastic_search_6
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
)

var IndexTest0Json interface{}

type IndexTestStruct struct{
	name string
	id int
	desc string
}

func main() {

	//创建一个空值的上下文，用于整合传递参数与适配运行状态
	ctx := context.Background()

	// 创建一个es6的客户端对象
	client, err := elastic.NewClient(
		/* 设置elastic服务实例ip
		 * DefaultURL = "http://127.0.0.1:9200"
		 */
		elastic.SetURL("http://10.3.138.104:9200/"),
		/* 监视时使用的协议，默认是http
		 * DefaultScheme = "http"
		 */
		elastic.SetScheme("http"),
	)

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
	}

	//1.client通过newclient创建默认开启，所以true
	fmt.Println(client.String(), client.IsRunning())

	//3.Ping一下es服务
	info, code, err := client.Ping("http://10.3.138.104:9200/").Do(ctx)
	if err != nil {
		panic(err)
	}
	//Elasticsearch returned with code 200 and version 6.2.2
	fmt.Printf("Elasticsearch 返回http返回值code %d 和版本 version %s\n", code, info.Version.Number)

	//4.获取es版本信息
	//可以直接curl localhost:9200 查看,也可以浏览器输入localhost:9200查看
	version, err := client.ElasticsearchVersion("http://10.3.138.104:9200/")
	if err != nil {
		panic(err)
	}
	//Elasticsearch version
	fmt.Printf("查看Elasticsearch版本version %s\n", version)


	//7.获取到指定id对象,Image To Imitate --Index("index_test_0").Type("type_test").Id("1")
	get1, err := client.Get().Index("index_test_0").Type("type_test").Id("1").Do(ctx)
	if err != nil {
		panic(err)
	}

	if get1.Found {
		//每次执行的version都不一样，因为前面的操作实际上重复put就会修改值，每次修改，es都会内部修改版本号
		fmt.Printf("Got document [%s] in version [%d] from index [%s], type [%s]\n", get1.Id, get1.Version, get1.Index, get1.Type)
		fmt.Println("尝试读取数据：")

		fmt.Printf("读取字段(json.RawMessage格式):%s",get1.Source)
        var raw = get1.Source
		err2 := json.Unmarshal(*raw, &IndexTest0Json)
		//json.Unmarshal() 函数将一个JSON对象解码到 空接口IndexTest0Json中，最终r将会是一个键值对的 map[string]interface{} 结构
		if err2 !=nil {
			panic(err2)
		}

		fmt.Println()
        fmt.Println("将得到的位置json按属性赋值给结构体对象: ")
		indexTestStruct := new(IndexTestStruct)
		unknowRow, ok := IndexTest0Json.(map[string]interface{})
		if ok {
			for k, v := range unknowRow {
				switch v2 := v.(type) {
				case string:
					if k=="name"{
						indexTestStruct.name = v2
					}
					if k=="desc"{
						indexTestStruct.desc = v2
					}
				case int:
					if k=="id"{
						indexTestStruct.id = v2
					}
				case float64:
					if k=="id"{
						indexTestStruct.id = int(v2)
					}
				case bool:
					fmt.Println(k, "is bool", v2)
				case []interface{}:
					fmt.Println(k, "is an array:")
					for i, iv := range v2 {
						fmt.Println(i, iv)
					}
				default:
					fmt.Println(k, "is another type not handle yet ")
					fmt.Printf("%s 's type is %T",k,v2)
				}
			}
		}

		fmt.Println(indexTestStruct)

	}

	//返回的值是分片结构，有总数-成功数-失败数
	_, err = client.Flush().Index("datacenter_book").Do(ctx)
	if err != nil {
		panic(err)
	}

	//关闭这个client，如后续还需要开启，执行client.Start()
	client.Stop()

}
