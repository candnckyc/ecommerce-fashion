package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"
	"ecommerce-backend/internal/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetProducts handles product listing with filters
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	query := &models.ProductListQuery{}

	// Parse query parameters
	queryParams := r.URL.Query()

	if categoryID := queryParams.Get("category"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			query.CategoryID = &id
		}
	}

	if brandID := queryParams.Get("brand"); brandID != "" {
		if id, err := strconv.Atoi(brandID); err == nil {
			query.BrandID = &id
		}
	}

	if minPrice := queryParams.Get("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query.MinPrice = &price
		}
	}

	if maxPrice := queryParams.Get("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query.MaxPrice = &price
		}
	}

	query.Size = queryParams.Get("size")
	query.Color = queryParams.Get("color")
	query.Search = queryParams.Get("search")

	if page := queryParams.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = p
		}
	}

	if limit := queryParams.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			query.Limit = l
		}
	}

	products, err := h.productService.GetProducts(query)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	utils.Success(w, products)
}

// GetProduct handles getting a single product
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}
	
	if product == nil {
		utils.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.Success(w, product)
}

// GetBrands handles getting all brands
func (h *ProductHandler) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.productService.GetBrands()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve brands")
		return
	}

	utils.Success(w, brands)
}

// GetCategories handles getting all categories
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productService.GetCategories()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	utils.Success(w, categories)
}
