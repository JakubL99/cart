// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/cart.proto

package cart

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Cart service

func NewCartEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Cart service

type CartService interface {
	Create(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error)
	GetCart(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error)
	DeleteCart(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error)
	DeleteProduct(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error)
}

type cartService struct {
	c    client.Client
	name string
}

func NewCartService(name string, c client.Client) CartService {
	return &cartService{
		c:    c,
		name: name,
	}
}

func (c *cartService) Create(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error) {
	req := c.c.NewRequest(c.name, "Cart.Create", in)
	out := new(SaveCart)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartService) GetCart(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error) {
	req := c.c.NewRequest(c.name, "Cart.GetCart", in)
	out := new(SaveCart)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartService) DeleteCart(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error) {
	req := c.c.NewRequest(c.name, "Cart.DeleteCart", in)
	out := new(SaveCart)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartService) DeleteProduct(ctx context.Context, in *ShoppingCart, opts ...client.CallOption) (*SaveCart, error) {
	req := c.c.NewRequest(c.name, "Cart.DeleteProduct", in)
	out := new(SaveCart)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Cart service

type CartHandler interface {
	Create(context.Context, *ShoppingCart, *SaveCart) error
	GetCart(context.Context, *ShoppingCart, *SaveCart) error
	DeleteCart(context.Context, *ShoppingCart, *SaveCart) error
	DeleteProduct(context.Context, *ShoppingCart, *SaveCart) error
}

func RegisterCartHandler(s server.Server, hdlr CartHandler, opts ...server.HandlerOption) error {
	type cart interface {
		Create(ctx context.Context, in *ShoppingCart, out *SaveCart) error
		GetCart(ctx context.Context, in *ShoppingCart, out *SaveCart) error
		DeleteCart(ctx context.Context, in *ShoppingCart, out *SaveCart) error
		DeleteProduct(ctx context.Context, in *ShoppingCart, out *SaveCart) error
	}
	type Cart struct {
		cart
	}
	h := &cartHandler{hdlr}
	return s.Handle(s.NewHandler(&Cart{h}, opts...))
}

type cartHandler struct {
	CartHandler
}

func (h *cartHandler) Create(ctx context.Context, in *ShoppingCart, out *SaveCart) error {
	return h.CartHandler.Create(ctx, in, out)
}

func (h *cartHandler) GetCart(ctx context.Context, in *ShoppingCart, out *SaveCart) error {
	return h.CartHandler.GetCart(ctx, in, out)
}

func (h *cartHandler) DeleteCart(ctx context.Context, in *ShoppingCart, out *SaveCart) error {
	return h.CartHandler.DeleteCart(ctx, in, out)
}

func (h *cartHandler) DeleteProduct(ctx context.Context, in *ShoppingCart, out *SaveCart) error {
	return h.CartHandler.DeleteProduct(ctx, in, out)
}
