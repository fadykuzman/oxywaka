package repositories

import (
	"testing"

	"github.com/muety/wakapi/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Count(t *testing.T) {
	repos := NewTestUserRepository(t)

	count, err := repos.Count()

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count, "empty database should have 0 users")
}

func TestUserRepository_FindOne_NonExistent(t *testing.T) {
	repo := NewTestUserRepository(t)

	_, err := repo.FindOne(models.User{ID: "nonexistent"})

	assert.Error(t, err, "should return error for non-existent user")
	assert.ErrorIs(t, err, ErrNotFound)
}
