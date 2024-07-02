package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       string             `bson:"price" json:"price"`
	Tag         string             `bson:"tag" json:"tag"`
	ImageUrl    string             `bson:"imageUrl" json:"imageUrl"`
}

type CartItem struct {
	Product  Product `bson:"product" json:"product"`
	Quantity int     `bson:"quantity" json:"quantity"`
}

type Cart struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Products []CartItem         `bson:"products" json:"products"`
}

// cart에 담아달리고 요청하는 상품의 이름들
type ProductNames struct {
	Tags []string `json:"tags"`
}
