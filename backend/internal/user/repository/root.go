package repository

import (
	"context"
	"fmt"

	"backend-golang/ent"
	entuser "backend-golang/ent/user"

	"golang.org/x/crypto/bcrypt"
)

// SyncRootUser applies configured credentials without replacing the user's ID.
// An empty email disables provisioning; config.Load validates the pair.
func SyncRootUser(ctx context.Context, client *ent.Client, email, password string) error {
	if email == "" {
		return nil
	}
	usr, err := client.User.Query().Where(entuser.Email(email)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("find root account: %w", err)
	}
	if usr != nil && bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)) == nil {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash root password: %w", err)
	}
	if usr == nil {
		_, err = client.User.Create().SetName("Root").SetEmail(email).SetPassword(string(hash)).Save(ctx)
		if err == nil {
			return nil
		}
		if !ent.IsConstraintError(err) {
			return fmt.Errorf("create root account: %w", err)
		}
		// Another backend instance may have created the same account on startup.
		usr, err = client.User.Query().Where(entuser.Email(email)).Only(ctx)
		if err != nil {
			return fmt.Errorf("find root account after concurrent creation: %w", err)
		}
	}
	if _, err := client.User.UpdateOneID(usr.ID).SetPassword(string(hash)).Save(ctx); err != nil {
		return fmt.Errorf("update root password: %w", err)
	}
	return nil
}
