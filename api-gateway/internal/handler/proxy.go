package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ServiceProxy handles proxying requests to the appropriate service
type ServiceProxy struct {
	services map[string]*httputil.ReverseProxy
	cfg      *config.Config
}

// NewServiceProxy creates a new service proxy
func NewServiceProxy(cfg *config.Config) *ServiceProxy {
	// Create reverse proxies for each service
	inventory, err := url.Parse(cfg.Services.Inventory.BaseURL)
	if err != nil {
		logrus.Fatalf("Invalid inventory service URL: %v", err)
	}

	order, err := url.Parse(cfg.Services.Order.BaseURL)
	if err != nil {
		logrus.Fatalf("Invalid order service URL: %v", err)
	}

	proxies := map[string]*httputil.ReverseProxy{
		"inventory": httputil.NewSingleHostReverseProxy(inventory),
		"order":     httputil.NewSingleHostReverseProxy(order),
	}

	return &ServiceProxy{
		services: proxies,
		cfg:      cfg,
	}
}

// ProxyInventory proxies requests to the inventory service
func (p *ServiceProxy) ProxyInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infof("Proxying request to inventory service: %s", c.Request.URL.Path)
		p.handleProxy(c, "inventory", p.cfg.Services.Inventory.BaseURL)
	}
}

// ProxyOrder proxies requests to the order service
func (p *ServiceProxy) ProxyOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infof("Proxying request to order service: %s", c.Request.URL.Path)
		p.handleProxy(c, "order", p.cfg.Services.Order.BaseURL)
	}
}

// handleProxy handles the actual proxying logic
func (p *ServiceProxy) handleProxy(c *gin.Context, service, baseURL string) {
	proxy := p.services[service]

	// Update the request URL path to match the target service
	// We need to remove the service prefix from the path
	path := c.Request.URL.Path
	prefix := "/" + service

	if strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
		if path == "" {
			path = "/"
		}
	}

	// Clone the request since we're modifying the URL
	outReq := new(http.Request)
	*outReq = *c.Request
	outReq.URL, _ = url.Parse(baseURL + path)
	if c.Request.URL.RawQuery != "" {
		outReq.URL.RawQuery = c.Request.URL.RawQuery
	}

	// Force the next handler to use our modified request
	c.Request = outReq

	// Let the reverse proxy do its job
	proxy.ServeHTTP(c.Writer, c.Request)
}

// RegisterRoutes registers the proxy routes
func RegisterProxyRoutes(router *gin.Engine, proxy *ServiceProxy) {
	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"service": "api-gateway",
		})
	})

	// Inventory service routes
	inventoryGroup := router.Group("/inventory")
	inventoryGroup.Any("/*path", proxy.ProxyInventory())

	// Order service routes
	orderGroup := router.Group("/order")
	orderGroup.Any("/*path", proxy.ProxyOrder())

	// Alternatively, we can also directly proxy requests to the root paths
	router.GET("/products", proxy.ProxyInventory())
	router.GET("/products/:id", proxy.ProxyInventory())
	router.POST("/products", proxy.ProxyInventory())
	router.PATCH("/products/:id", proxy.ProxyInventory())
	router.DELETE("/products/:id", proxy.ProxyInventory())

	router.GET("/categories", proxy.ProxyInventory())
	router.GET("/categories/:id", proxy.ProxyInventory())
	router.POST("/categories", proxy.ProxyInventory())
	router.PATCH("/categories/:id", proxy.ProxyInventory())
	router.DELETE("/categories/:id", proxy.ProxyInventory())

	router.GET("/orders", proxy.ProxyOrder())
	router.GET("/orders/:id", proxy.ProxyOrder())
	router.POST("/orders", proxy.ProxyOrder())
	router.PATCH("/orders/:id", proxy.ProxyOrder())
}
