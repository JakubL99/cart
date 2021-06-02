package handler

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	productProto "cart/product"
	pb "cart/proto"
)

type Cart struct {
	ID       primitive.ObjectID `bson:"_id"`
	IdUser   string             `json:"idUser"`
	Products Products           `json:"products"`
}

type SaveCart struct {
	IdUser   string   `json:"idUser"`
	Products Products `json:"products"`
}

type Product struct {
	IdProduct string `json:"idProduct"`
	Name      string `json:"name"`
	Price     string `json:"price"`
}

type idUser struct {
	ID string `bson:"id"`
}

type idProduct struct {
	ID string `bson:"id"`
}

type Products []Product

type Repository interface {
	Create(ctx context.Context, cart *SaveCart) error
	FindCart(ctx context.Context, id *idUser) *Cart
	FindOne(ctx context.Context, id *idUser) (*Cart, error)
	UpdateCart(ctx context.Context, id *idUser, products []Product) error
	DeleteOne(ctx context.Context, id *idUser) (*Cart, error)
}

type MongoRepository struct {
	Collection *mongo.Collection
}

type Handler struct {
	Repository
	ProductClient productProto.ProductService
}

func MarshalIdProduct(id *pb.ShoppingCart) *productProto.GetRequest {
	return &productProto.GetRequest{
		Id: id.IdProduct,
	}
}

func MarshalIdP(id *pb.ShoppingCart) *idProduct {
	return &idProduct{
		ID: id.IdProduct,
	}
}

func MarshalIdUser(id *pb.ShoppingCart) *idUser {
	return &idUser{
		ID: id.IdUser,
	}
}

func MarshalProduct(product *productProto.GetResponse) Product {
	return Product{
		IdProduct: product.Product.Id,
		Name:      product.Product.Name,
		Price:     product.Product.Price,
	}
}

func UnmarshalProduct(product Product) *pb.Product {
	return &pb.Product{
		IdProduct: product.IdProduct,
		Name:      product.Name,
		Price:     product.Price,
	}
}

func UnmarshalCart(cart *Cart) *pb.SaveCart {
	products := make([]*pb.Product, 0)
	for _, product := range cart.Products {
		products = append(products, UnmarshalProduct(product))
	}
	obj_id := primitive.ObjectID.Hex(cart.ID)
	return &pb.SaveCart{
		Id:       obj_id,
		IdUser:   cart.IdUser,
		Products: products,
	}
}

func removeElement(s []Product, i int) ([]Product, error) {

	if i >= len(s) || i < 0 {
		return nil, fmt.Errorf("Index is out of range. Index is %d with slice length %d", i, len(s))
	}

	s[i] = s[len(s)-1]

	return s[:len(s)-1], nil
}

func (repo *MongoRepository) Create(ctx context.Context, cart *SaveCart) error {
	_, err := repo.Collection.InsertOne(ctx, cart)
	return err
}

func (repo *MongoRepository) FindCart(ctx context.Context, id *idUser) *Cart {
	cur := repo.Collection.FindOne(ctx, bson.M{"iduser": id.ID})
	var cart *Cart
	if err := cur.Decode(&cart); err != nil {
		return nil
	}

	return cart
}

func (repo *MongoRepository) FindOne(ctx context.Context, id *idUser) (*Cart, error) {
	cur := repo.Collection.FindOne(ctx, bson.M{"iduser": id.ID})
	var cart *Cart
	if err := cur.Decode(&cart); err != nil {
		return nil, err
	}

	return cart, cur.Err()
}

func (repo *MongoRepository) UpdateCart(ctx context.Context, id *idUser, products []Product) error {
	result := repo.Collection.FindOneAndUpdate(
		ctx,
		bson.M{"iduser": id.ID},
		bson.M{"$set": bson.M{"products": products}},
	)
	if result.Err() != nil {
		logger.Error("Error Update Cart:  ", result.Err())
	}
	logger.Info("Update cart:  ", result)
	return result.Err()
}
func (repo *MongoRepository) DeleteOne(ctx context.Context, id *idUser) (*Cart, error) {
	cur := repo.Collection.FindOne(ctx, bson.M{"iduser": id.ID})
	var cart *Cart
	if err := cur.Decode(&cart); err != nil {
		logger.Error("Error Decode ", err)
		return nil, err
	}
	resultDelete, err := repo.Collection.DeleteOne(ctx, bson.M{"iduser": id.ID})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	logger.Info("resultDelete: ", resultDelete.DeletedCount)
	return cart, cur.Err()
}

func (h *Handler) Create(ctx context.Context, req *pb.ShoppingCart, rsp *pb.SaveCart) error {
	productResponse, err := h.ProductClient.GetProduct(ctx, MarshalIdProduct(req))
	if err != nil {
		logger.Error("Error Get Product: ", err)
	}
	logger.Info("productResponse =  ", productResponse)
	var User *idUser
	User = MarshalIdUser(req)

	var Pro Product
	Pro = MarshalProduct(productResponse)

	var products []Product
	UserCart := h.FindCart(ctx, User)
	if UserCart == nil {
		products = append(products, Pro)

		var CreateCart SaveCart
		CreateCart.IdUser = User.ID
		CreateCart.Products = products

		if err := h.Repository.Create(ctx, &CreateCart); err != nil {
			return err
		}
		cart, err := h.Repository.FindOne(ctx, User)
		if err != nil {
			logger.Error("Error function FindOne", err)
		}
		unc := UnmarshalCart(cart)
		rsp = unc
		products = nil
	} else {
		products = UserCart.Products
		products = append(products, Pro)
		if err := h.Repository.UpdateCart(ctx, User, products); err != nil {
			logger.Info("Error function UpdateCart:  ", err)
		}
		cart, err := h.Repository.FindOne(ctx, User)
		if err != nil {
			logger.Error("Error function FindOne", err)
		}
		unc := UnmarshalCart(cart)
		rsp.Id = unc.Id
		rsp.IdUser = unc.IdUser
		rsp.Products = unc.Products
		products = nil
	}

	return nil
}

func (h *Handler) GetCart(ctx context.Context, req *pb.ShoppingCart, rsp *pb.SaveCart) error {
	var User *idUser
	User = MarshalIdUser(req)
	cart, err := h.Repository.FindOne(ctx, User)
	if err != nil {
		logger.Error("Error function FindOne", err)
	}
	unc := UnmarshalCart(cart)
	rsp.Id = unc.Id
	rsp.IdUser = unc.IdUser
	rsp.Products = unc.Products
	return nil
}

func (h *Handler) DeleteCart(ctx context.Context, req *pb.ShoppingCart, rsp *pb.SaveCart) error {
	var User *idUser
	User = MarshalIdUser(req)
	cart, err := h.Repository.DeleteOne(ctx, User)
	if err != nil {
		logger.Error("Error function FindOne", err)
	}
	unc := UnmarshalCart(cart)
	rsp.Id = unc.Id
	rsp.IdUser = unc.IdUser
	rsp.Products = unc.Products
	return nil
}

func (h *Handler) DeleteProduct(ctx context.Context, req *pb.ShoppingCart, rsp *pb.SaveCart) error {
	var User *idUser
	User = MarshalIdUser(req)
	cart, err := h.Repository.FindOne(ctx, User)
	if err != nil {
		logger.Error("Error function FindOne", err)
	}
	var products []Product
	products = cart.Products
	id := MarshalIdP(req)
	for i, v := range products {
		if v.IdProduct == id.ID {
			p, err := removeElement(products, i)
			if err != nil {
				logger.Error("Error remove product from cart:  ", err)
			}
			if err := h.Repository.UpdateCart(ctx, User, p); err != nil {
				logger.Info("Error function UpdateCart:  ", err)
			}
		}
	}
	cart, err = h.Repository.FindOne(ctx, User)
	if err != nil {
		logger.Error("Error function FindOne", err)
	}
	unc := UnmarshalCart(cart)
	rsp.Id = unc.Id
	rsp.IdUser = unc.IdUser
	rsp.Products = unc.Products
	return nil
}
