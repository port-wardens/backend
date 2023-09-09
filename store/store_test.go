package store_test

import (
	"github.com/jayendramadaram/port-wardens/model"
	"github.com/jayendramadaram/port-wardens/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var (
		s store.Store
	)

	BeforeSuite(func() {
		dns := "host=localhost user=kakkoii4cat password=root dbname=orderbook port=5432 sslmode=disable TimeZone=Asia/Kolkata"
		db, err := model.NewDB(dns)
		if err != nil {
			panic(err)
		}
		s = store.NewStore(db)
	})

	Describe("HealthCheck", func() {
		It("should return nil", func() {
			Expect(s.HealthCheck()).To(Succeed())
		})
	})

})
