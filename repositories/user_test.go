package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Count(t *testing.T) {
	repos := NewTestUserRepository(t)

	count, err := repos.Count()

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count, "empty database should have 0 users")
}

// func TestUserRepository_InsertOrGet_NewUser(t *testing.T) {
// 	repo := NewTestUserRepository(t)
//
// 	user := &models.User{
// 		ID:       "testuser",
// 		Email:    "test@example.com",
// 		Password: "hashedpassword",
// 	}
//
// 	result, created, err := repo.InsertOrGet(user)
//
// 	assert.NoError(t, err)
// 	assert.True(t, created, "user should be created")
// 	assert.Equal(t, "testuser", result.ID)
// 	assert.Equal(t, "test@example.com", result.Email)
// }
