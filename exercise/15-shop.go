package main

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type product struct {
	Name  string
	Price float64
	Count int
}

type Store struct {
	products map[string]*product
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]*product),
	}
}

func (s *Store) AddProduct(name string, price float64, count int) error {
	if s.products == nil {
		s.products = make(map[string]*product)
	}

	key := strings.ToLower(name)

	if _, ok := s.products[key]; ok {
		return errors.New(name + " already exists")
	}

	if price <= 0 {
		return errors.New("price should be positive")
	}

	if count <= 0 {
		return errors.New("count should be positive")
	}

	s.products[key] = &product{
		Name:  name,
		Price: price,
		Count: count,
	}

	return nil
}

func (s *Store) GetProductCount(name string) (int, error) {
	key := strings.ToLower(name)
	p, ok := s.products[key]
	if !ok {
		return 0, errors.New("invalid product name")
	}
	return p.Count, nil
}

func (s *Store) GetProductPrice(name string) (float64, error) {
	key := strings.ToLower(name)
	p, ok := s.products[key]
	if !ok {
		return 0, errors.New("invalid product name")
	}
	return p.Price, nil
}

func (s *Store) Order(name string, count int) error {
	if count <= 0 {
		return errors.New("count should be positive")
	}

	key := strings.ToLower(name)
	p, ok := s.products[key]
	if !ok {
		return errors.New("invalid product name")
	}

	if p.Count == 0 {
		return errors.New("there is no " + name + " in the store")
	}

	if count > p.Count {
		return errors.New("not enough " + name + " in the store. there are " +
			strconv.Itoa(p.Count) + " left")
	}

	p.Count -= count
	return nil
}

func (s *Store) ProductsList() ([]string, error) {
	var result []string

	for _, p := range s.products {
		if p.Count > 0 {
			result = append(result, strings.ToLower(p.Name))
		}
	}

	if len(result) == 0 {
		return nil, errors.New("store is empty")
	}

	sort.Strings(result)
	return result, nil
}
