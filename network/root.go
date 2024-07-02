package network

import (
	"context"
	"cpk_mall/types"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string

func MakeRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/get/products", GetAllProduct)
	r.POST("/cart/input", PostItmes)
	// uri 가져오기
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri = os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	return r
}

// 모든 상품을 가져오는 핸들러
func GetAllProduct(c *gin.Context) {

	// mongoDB 서버 연결
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// 끝나기 전에 클라이언트 종료
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// collection 가져오기
	coll := client.Database("cpk-mall").Collection("products")

	// 모든 문서 찾기
	ctx := context.Background()
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// cursor에 collection의 모든 정보들이 들어감
	// ctx는 단지 흐름을 제어하기 위한 용도 (ex: ctx에 timeout 5초 설정하면 5초 후에 false를 반환)
	var data []bson.M = make([]bson.M, 0, 100)
	for cursor.Next(ctx) {
		var result bson.M // bson.M : map[string]interface{}의 별칭
		err := cursor.Decode(&result)
		data = append(data, result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	// 읽는 도중 Err가 있었는지 확인
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, data)
}

func PostItmes(c *gin.Context) {
	// search product with name
	var productNames types.ProductNames
	// get data from request body
	// JSON 바인딩 => 언마샬링해서 productNames에 저장
	if err := c.BindJSON(&productNames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DB에서 상품 조회

	// mongoDB 서버 연결
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// collection 가져오기
	coll := client.Database("cpk-mall").Collection("products")

	var searchedProduct [][]types.Product
	var query bson.M

	// tag로 상품 조회
	for _, v := range productNames.Tags {
		query = bson.M{"tag": v}
		// 쿼리 실행
		cursor, err := coll.Find(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(context.Background())

		var products []types.Product
		if err = cursor.All(context.Background(), &products); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		searchedProduct = append(searchedProduct, products)
	}

	c.JSON(http.StatusOK, searchedProduct)

	// 조회된 상품들 가중치 매겨서 내림차순 정렬

	// 맨 앞을 cart에 담음

	// 잘 담았다고 알려줌
}
