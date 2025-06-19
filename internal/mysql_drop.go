package internal

import (
	"fmt"
	"stackroost/internal/logger"
)

func DropMySQLUserAndDatabase(username string) error {
	logger.Info(fmt.Sprintf("Dropping MySQL user and database: %s", username))

	sql := fmt.Sprintf(`
		DROP DATABASE IF EXISTS %s;
		DROP USER IF EXISTS '%s'@'localhost';
		FLUSH PRIVILEGES;
	`, username, username)

	cmd := []string{"-e", sql}
	if err := RunCommand("sudo", append([]string{"mysql"}, cmd...)...); err != nil {
		return fmt.Errorf("failed to drop MySQL user or database: %v", err)
	}
	return nil
}
