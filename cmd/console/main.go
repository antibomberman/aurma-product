package main

import (
	"aurma_product/internal/di"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Printf("Failed to initialize container: %v", err)
		os.Exit(1)
	}
	defer container.Close()

	console := &cobra.Command{Use: "app", Short: "Aurma Product Application"}

	console.AddCommand(&cobra.Command{
		Use:   "elastic-index-create",
		Short: "create Elasticsearch index for products",
		Run: func(cmd *cobra.Command, args []string) {
			err := container.Elastic.ProductCreateIndex()
			if err != nil {
				fmt.Println(err)
			}
		},
	}, &cobra.Command{
		Use:   "elastic-products",
		Short: "add all products to Elasticsearch",
		Run: func(cmd *cobra.Command, args []string) {
			err := container.ProductService.SetAllProductToElastic()
			if err != nil {
				log.Printf("Error setting products to Elastic: %v", err)
				return
			}
		},
	})

	err = console.Execute()
	if err != nil {
		log.Printf("Error executing command: %v", err)
		return
	}
}
