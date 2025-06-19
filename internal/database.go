package internal

import (
	"fmt"
	"stackroost/internal/logger"
)
func CreateMySQLUserAndDatabase(username, password string) error {
	logger.Info(fmt.Sprintf("Creating MySQL user and database for '%s'", username))

	createUserCmd := fmt.Sprintf(`CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s';`, username, password)
	createDBCmd := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s;`, username)
	grantCmd := fmt.Sprintf(`GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost';`, username, username)
	flushCmd := `FLUSH PRIVILEGES;`
	fullSQL := fmt.Sprintf("%s %s %s %s", createUserCmd, createDBCmd, grantCmd, flushCmd)
	mysqlCmd := []string{"-e", fullSQL}
	if err := RunCommand("sudo", append([]string{"mysql"}, mysqlCmd...)...); err != nil {
		return fmt.Errorf("failed to create MySQL user or database: %v", err)
	}

	logger.Success(fmt.Sprintf("MySQL user and database '%s' created", username))
	return nil
}
