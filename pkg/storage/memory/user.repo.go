package memory

import (
	"github.com/slim-crown/issue-1-REST/pkg/domain/user"
	"time"
)


// UserRepository ...
type UserRepository struct {
	cache         map[string]user.User
	secondaryRepo *user.Repository
	allRepos      *map[string]interface{}
}

// NewUserRepository returns a new in memory cache implementation of user.Repository.
// The database implementation of user.Repository must be passed as the first argument
// since to simplify logic, cache repos wrap the database repos.
// A map of all the other cache based implementations of the Repository interfaces
// found in the different services of the project must be passed as a second argument as
// the Repository might make user of them to fetch objects instead of implementing redundant logic.
func NewUserRepository(dbRepo *user.Repository, allRepos *map[string]interface{}) *UserRepository {
	return &UserRepository{make(map[string]user.User, 100), dbRepo, allRepos}
}

// AddUser takes in a user.User struct and persists it.
// Returns an error if the DB repository implementation returns an error.
func (repo *UserRepository) AddUser(u *user.User) error {
	return (*repo.secondaryRepo).AddUser(u)
}

// GetUser retrives a user.User based on the username passed.
func (repo *UserRepository) GetUser(username string) (*user.User, error) {
	if _, ok := repo.cache[username]; ok == false {
		err := repo.cacheUser(username)
		if err != nil {
			return nil, err
		}
	}
	u := repo.cache[username]
	return &u, nil

}

// cacheUser is just a helper function
func (repo *UserRepository) cacheUser(username string) error {
	u, err := (*repo.secondaryRepo).GetUser(username)
	if err != nil {
		return err
	}
	repo.cache[username] = *u
	return nil
}

// UpdateUser updates a user based on the passed user.User struct.
func (repo *UserRepository) UpdateUser(username string, u *user.User) error {
	err := (*repo.secondaryRepo).UpdateUser(username, u)
	if err == nil {
		// If updating in the DB repo is successful, it updates its cache by getting
		// the new user.User and converting it into a cache able format.
		if u.Username != "" {
			// if the username is changed, use the new username from the struct to update the cache
			err = repo.cacheUser(u.Username)
			if err != nil {
				return err
			}
		} else {
			delete(repo.cache, u.Username)
			err = repo.cacheUser(username)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// DeleteUser deletes a user based on the passed in username.
func (repo *UserRepository) DeleteUser(username string) error {
	err := (*repo.secondaryRepo).DeleteUser(username)
	if err == nil {
		// If deletion is successful, it also tries to delete the user from its cache.
		delete(repo.cache, username)
	}
	return err
}

// SearchUser calls the DB repo SearchUser function.
// It also caches all the users returned by the result.
func (repo *UserRepository) SearchUser(pattern, sortBy, sortOrder string, limit, offset int) ([]*user.User, error) {
	// TODO: memory.UserRepository.SerarchUser method
	result, err := (*repo.secondaryRepo).SearchUser(pattern, sortBy, sortOrder, limit, offset)
	if err == nil {
		for _, u := range result {
			uTemp := *u
			repo.cache[u.Username] = uTemp
			u.Email = ""
			u.BookmarkedPosts = make(map[int]time.Time)
		}
	}
	return result, err
}

// PassHashIsCorrect calls the DB repo PassHashIsCorrect function. checks the given pass hash against the pass hash found in the database for the username.
func (repo *UserRepository) PassHashIsCorrect(username, passHash string) bool {
	return (*repo.secondaryRepo).PassHashIsCorrect(username, passHash)
}

// BookmarkPost calls the DB repo BookmarkPost function.
func (repo *UserRepository) BookmarkPost(username string, postID int) error {
	err := (*repo.secondaryRepo).BookmarkPost(username, postID)
	if err == nil {
		err = repo.cacheUser(username)
		if err != nil {
			return err
		}
	}
	return err
}

// DeleteBookmark calls the DB repo BookmarkPost function.
func (repo *UserRepository) DeleteBookmark(username string, postID int) error {
	err := (*repo.secondaryRepo).DeleteBookmark(username, postID)
	if err == nil {
		err = repo.cacheUser(username)
		if err != nil {
			return err
		}
	}
	return err
}

// UsernameOccupied calls the DB repo UsernameOccupied function.
func (repo *UserRepository) UsernameOccupied(username string) (bool, error) {
	return (*repo.secondaryRepo).UsernameOccupied(username)
}

// EmailOccupied calls the DB repo EmailOccupied function.
func (repo *UserRepository) EmailOccupied(email string) (bool, error) {
	return (*repo.secondaryRepo).EmailOccupied(email)
}
