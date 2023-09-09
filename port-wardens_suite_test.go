package main_test

import (
	"fmt"
	"testing"

	"github.com/jayendramadaram/port-wardens/model"
	"github.com/jayendramadaram/port-wardens/testclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOrderbook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orderbook Suite")
}

var _ = Describe("Orderbook", func() {
	var (
		c testclient.Client
	)

	BeforeSuite(func() {
		c = testclient.NewPortWardenClient("http://localhost:8080")
	})

	It("Signup should return 200", func() {
		// Skip("Skipping Signup Cuz user already exists")
		response := c.SignUp(model.CreateUser{
			Email:    "panther@catalog.fi",
			Password: "password",
			Username: "panther",
		})
		Expect(response).NotTo(HaveOccurred())
	})

	It("Login should return 200", func() {
		response, err := c.Login(model.LoginUser{
			Email:    "panther@gmail.com",
			Password: "password",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Token).NotTo(BeEmpty())
		fmt.Println("Token: ", response.Token)
	})

})
