package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Size struct {
	OnStock  bool   `json:"onStock"`
	SizeName string `json:"sizeName"`
}

type Color struct {
	Hex  string `json:"hex"`
	Name string `json:"name"`
}

type Stock struct {
	Quantity  int  `json:"quantity"`
	IsInStock bool `json:"isInStock"`
}

type Product struct {
	ID              string  `json:"_id"`
	Name            string  `json:"name"`
	Brand           string  `json:"brand"`
	Price           float64 `json:"price"`
	Currency        string  `json:"currency"`
	PriceInRubles   float64 `json:"priceInRubles"`
	DiscountedPrice *float64 `json:"discountedPrice"`
	Description     string  `json:"description"`
	Images          []string `json:"images"`
	Sizes           []Size  `json:"sizes"`
	Colors          []Color `json:"colors"`
	ProductURL      string  `json:"productUrl"`
	Store           string  `json:"store"`
	Category        string  `json:"category"`
	ProcessedAt     string  `json:"processedAt"`
	IsActive        bool    `json:"isActive"`
	StockStatus     string  `json:"stockStatus"`
	Stock           Stock   `json:"stock"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	stores := []string{"zara", "nike", "adidas"}
	brands := map[string]string{
		"zara":   "Zara",
		"nike":   "Nike", 
		"adidas": "Adidas",
	}
	
	categories := []string{"man-shirts", "woman-dresses", "man-pants", "woman-tops", "shoes", "accessories"}
	colors := []Color{
		{Hex: "#FF0000", Name: "Kırmızı"},
		{Hex: "#00FF00", Name: "Yeşil"},
		{Hex: "#0000FF", Name: "Mavi"},
		{Hex: "#FFFFFF", Name: "Beyaz"},
		{Hex: "#000000", Name: "Siyah"},
		{Hex: "#FFFF00", Name: "Sarı"},
		{Hex: "#FF00FF", Name: "Pembe"},
		{Hex: "#00FFFF", Name: "Turkuaz"},
	}
	
	sizes := []string{"XS", "S", "M", "L", "XL", "XXL"}
	stockStatuses := []string{"in_stock", "low_stock", "out_of_stock"}
	
	for _, store := range stores {
		var products []Product
		
		fmt.Printf("%s mağazası için 5000 ürün oluşturuluyor...\n", store)
		
		for i := 1; i <= 5000; i++ {
			// Random product data
			price := 50.0 + rand.Float64()*450.0 // 50-500 TL arası
			quantity := rand.Intn(100) + 1
			isInStock := quantity > 0
			stockStatus := stockStatuses[rand.Intn(len(stockStatuses))]
			if !isInStock {
				stockStatus = "out_of_stock"
			}
			
			// Random sizes
			productSizes := make([]Size, 0)
			numSizes := rand.Intn(4) + 1 // 1-4 beden
			for j := 0; j < numSizes; j++ {
				productSizes = append(productSizes, Size{
					OnStock:  rand.Intn(2) == 1,
					SizeName: sizes[rand.Intn(len(sizes))],
				})
			}
			
			// Random colors
			productColors := make([]Color, 0)
			numColors := rand.Intn(3) + 1 // 1-3 renk
			for j := 0; j < numColors; j++ {
				productColors = append(productColors, colors[rand.Intn(len(colors))])
			}
			
			product := Product{
				ID:              fmt.Sprintf("%s_product_%d_%d", store, i, rand.Intn(1000000)),
				Name:            fmt.Sprintf("%s Ürün %d - Premium Kalite", brands[store], i),
				Brand:           brands[store],
				Price:           price,
				Currency:        "TRY",
				PriceInRubles:   price * 15.5, // Approximate conversion
				DiscountedPrice: nil,
				Description:     fmt.Sprintf("Bu %s markasının premium kalitesinde %d numaralı ürünüdür. Yüksek kalite malzemelerden üretilmiştir.", brands[store], i),
				Images:          []string{fmt.Sprintf("https://example.com/images/%s_product_%d_image_1.jpg", store, i)},
				Sizes:           productSizes,
				Colors:          productColors,
				ProductURL:      fmt.Sprintf("https://%s.com/product-%d", store, i),
				Store:           store,
				Category:        categories[rand.Intn(len(categories))],
				ProcessedAt:     "21:37:51",
				IsActive:        true,
				StockStatus:     stockStatus,
				Stock: Stock{
					Quantity:  quantity,
					IsInStock: isInStock,
				},
			}
			
			products = append(products, product)
		}
		
		// Save to file
		filename := fmt.Sprintf("products_%s_5000.json", store)
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Dosya oluşturma hatası: %v\n", err)
			continue
		}
		defer file.Close()
		
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(products); err != nil {
			fmt.Printf("JSON yazma hatası: %v\n", err)
			continue
		}
		
		fmt.Printf("%s mağazası için %d ürün %s dosyasına kaydedildi\n", store, len(products), filename)
	}
	
	fmt.Println("Tüm mağazalar için ürün verileri oluşturuldu!")
}