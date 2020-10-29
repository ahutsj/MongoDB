package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	//"go.mongodb.org/mongo-driver/bson"
)

type AHUTSJ struct {
	Name string
	Age  int
}

func main() {
	//设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	//连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB！")

	//指定获取要操作的数据集
	collecitons := client.Database("tsj").Collection("AHUTSJ")

	s1 := AHUTSJ{"张三",12}
	s2 := AHUTSJ{"李四",13}
	s3 := AHUTSJ{"王五",14}

	//插入文档
	//插入文档(一条文档)
	insertResult,err := collecitons.InsertOne(context.TODO(),s1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document:",insertResult.InsertedID)

	//插入文档(多条文档)
	Multi := []interface{}{s2,s3}
	insertMultiResult,err := collecitons.InsertMany(context.TODO(),Multi)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents:",insertMultiResult.InsertedIDs)


	//更新文档
	//updateOne()方法允许你更新单个文档。它需要一个筛选器文档来匹配数据库中的文档，并需要一个更新文档来描述更新操作
	//可以使用bson.D类型来构建筛选文档和更新文档
	filter := bson.D{{"name","李四"}}
	update := bson.D{
		{"$inc",bson.D{
			{"age",20}, //age 增加了(increase)20
		}},
	}
	updateResult,err := collecitons.UpdateOne(context.TODO(),filter,update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and update %v documents!\n",updateResult.MatchedCount,updateResult.ModifiedCount)


	//查找文档
	//要找到一个文档，需要一个filter文档，以及一个指向可以将结果解码为其值得指针
	//要查找单个文档，使用collection.FindOne()，这个方法返回一个可以解码为值得结果
	var result AHUTSJ
	err = collecitons.FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document:%+v\n",result)

	//查找多个文档，使用collection.Find(),此方法返回一个游标，游标提供了一个文档流，可以通过它一次迭代和解码一个文档。
	//当游标用完后，应该关闭游标，下面的示例将使用options包设置一个限制以便只返回两个文档
	// 将选项传递给Find()
	findOptions := options.Find()
	findOptions.SetLimit(2)
	var results []*AHUTSJ //定义一个切片来存储查询结果
	cur,err := collecitons.Find(context.TODO(),bson.D{{}},findOptions) //把bson.D{{}}作为一个filter来匹配所有文档
	if err != nil {
		log.Fatal(err)
	}
	//查找多个文档返回一个光标
	//遍历游标允许一次解码一个文档
	for cur.Next(context.TODO()) {
		//创建一个值，将单个文档解码为该值
		var elem AHUTSJ
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results,&elem)
	}
	if err := cur.Err();err != nil {
		log.Fatal(err)
	}
	//完成后关闭游标
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents(array of pointers):%#v\n",results)


	//删除文档
	//可以使用collection.DeleteOne()或collection.DeleteMany()删除文档。
	//如果传递bson.D{{}}作为过滤器参数，它将匹配数据集中的所有文档。还可以使用collection.drop()删除整个数据集
	deleteResult1, err := collecitons.DeleteOne(context.TODO(), bson.D{{"name","小黄"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)
	// 删除所有
	deleteResult2, err := collecitons.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)


	//断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connection to MongoDB closed!")

}

