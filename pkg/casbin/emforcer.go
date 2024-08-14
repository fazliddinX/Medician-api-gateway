package casbin

import (
	config2 "api-gateway/pkg/config"
	xormadapter "github.com/casbin/xorm-adapter/v2"

	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
	"log/slog"
)

func CasbinEnforcer(logger *slog.Logger) (*casbin.Enforcer, error) {
	config := config2.Load()
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_NAME, config.DB_PASSWORD)

	adapter, err := xormadapter.NewAdapter("postgres", conn)
	if err != nil {
		log.Println("error creating Casbin adapter", "error", err.Error())
		logger.Error("Error creating Casbin adapter", "error", err.Error())
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("pkg/casbin/model.conf", adapter)
	if err != nil {
		logger.Error("Error creating Casbin enforcer", "error", err.Error())
		log.Println("error creating Casbin enforcer", "error", err.Error())
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Println("error loading Casbin policy", "error", err.Error())
		logger.Error("Error loading Casbin policy", "error", err.Error())
		return nil, err
	}

	policies := [][]string{
		{"patient", "/api/user", "GET"},
		{"patient", "/api/user", "PUT"},
		{"patient", "/api/user", "DELETE"},
		{"patient", "/api/medical-records", "POST"},
		{"patient", "/api/medical-records/:id", "GET"},
		{"patient", "/api/medical-records/:id", "PUT"},
		{"patient", "/api/medical-records/:id", "DELETE"},
		{"patient", "/api/lifestyle", "POST"},
		{"patient", "/api/lifestyle/:id", "GET"},
		{"patient", "/api/lifestyle/:id", "PUT"},
		{"patient", "/api/lifestyle/:id", "DELETE"},
		{"patient", "/api/health-recommendations/:id", "GET"},
		{"patient", "/api/health-recommendations", "GET"},

		{"doctor", "/api/medical-records/:id", "GET"},
		{"doctor", "/api/wearable-data", "POST"},
		{"doctor", "/api/wearable-data/:id", "GET"},
		{"doctor", "/api/wearable-data/:id", "PUT"},
		{"doctor", "/api/wearable-data/:id", "DELETE"},
		{"doctor", "/api/lifestyle/:id", "GET"},
		{"doctor", "/api/health-recommendations/:id", "POST"},
		{"doctor", "/api/health-recommendations/:id", "GET"},
		{"doctor", "/api/health-recommendations", "GET"},

		{"admin", "/api/user/all", "GET"},
		{"admin", "/api/user/create", "POST"},
		{"admin", "/api/medical-records", "GET"},
		{"admin", "/api/wearable-data", "GET"},
		{"admin", "/api/lifestyle", "GET"},
		{"admin", "/api/health-recommendations/:id", "GET"},
		{"admin", "/api/health-recommendations", "GET"},
		{"admin", "/api/monitoring/realtime", "DELETE"},
		{"admin", "/api/monitoring/daily-summary", "GET"},
		{"admin", "/api/monitoring/weekly-summary", "GET"},
	}

	_, err = enforcer.AddPolicies(policies)
	if err != nil {
		log.Println("error adding Casbin policy", "error", err.Error())
		logger.Error("Error adding Casbin policy", "error", err.Error())
		return nil, err
	}

	err = enforcer.SavePolicy()
	if err != nil {
		log.Println("Error saving Casbin policy", "error", err.Error())
		logger.Error("Error saving Casbin policy", "error", err.Error())
		return nil, err
	}
	return enforcer, nil
}
