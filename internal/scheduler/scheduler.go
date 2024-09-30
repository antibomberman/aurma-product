package scheduler

import (
	"aurma_product/internal/di"
	"github.com/robfig/cron/v3"
	"log"
)

func Run(container *di.Container) {
	c := cron.New()
	// Примеры использования:
	//@monthly
	//@weekly
	//@daily
	//@hourly
	//@every 1h30m
	//@every 30m
	//@every 30s
	//------------------------------------
	//"* 0 22 * * *" - это scheduler-выражение, которое означает "каждый день в 22:00". Вот расшифровка:
	//*: секунда (0-59)
	//0: минута (0-59)
	//22: час (0-23)
	//*: день месяца (1-31)
	//*: месяц (1-12)
	//*: день недели (0-7, где 0 и 7 - воскресенье)

	c.AddFunc("@every 3m", func() {
		products, err := container.ProductService.UpdatedProductPharmacies()
		if err != nil {
			log.Printf("Error: Error getting updated product pharmacies: %v", err)
			return
		}
		if len(products) > 0 {
			if err := container.Elastic.ProductAddDocument(products); err != nil {
				log.Printf("Error: Error indexing updated product pharmacies: %v", err)
			}
		}
	})

	c.Start()

	// Пустой select {} блокирует выполнение текущей горутины бесконечно
	select {}
}
