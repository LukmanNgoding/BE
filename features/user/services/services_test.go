package services

import (
	"errors"
	"testing"

	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/config"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelete(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("success", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(nil).Once()
		srv := New(repo)
		input := uint(1)
		err := srv.Delete(input)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})
	t.Run("Failed Delete", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("no data")).Once()
		srv := New(repo)
		input := uint(7)
		err := srv.Delete(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "no data", "error message doesn't match")
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Update", func(t *testing.T) {
		repo.On("Update", mock.Anything).Return(domain.Core{ID: uint(1), Username: "same", Email: "skjsa"}, nil).Once()
		srv := New(repo)
		input := domain.Core{ID: 1, Username: "same", Email: "skjsa", Password: "same"}
		res, err := srv.UpdateProfile(input)
		assert.Nil(t, err)
		assert.NotEmpty(t, res.ID, "Seharusnya ada ID user yang diupdate")
		assert.NotEqual(t, res.Password, input.Password, "Password tidak terenkripsi")
		assert.Equal(t, input.Username, res.Username, "Nama user harus sesuai")
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update", func(t *testing.T) {
		repo.On("Update", mock.Anything).Return(domain.Core{}, errors.New("rejected from database")).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Email: "skjsa", Password: "same"}
		res, err := srv.UpdateProfile(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "rejected from database", "Pesan error tidak sesuai")
		assert.Equal(t, uint(0), res.ID, "ID seharusnya 0 karena gagal input data")
		repo.AssertExpectations(t)
	})
}

func TestAddUserr(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Insert", func(t *testing.T) {
		repo.On("Insert", mock.Anything).Return(domain.Core{ID: uint(1), Username: "same", Email: "skjsa"}, nil).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Email: "skjsa", Password: "same"}
		res, err := srv.AddUser(input)
		assert.Nil(t, err)
		assert.NotEmpty(t, res.ID, "Seharusnya ada ID yang berhasil dibuat")
		assert.NotEqual(t, res.Password, input.Password, "Password tidak terenkripsi")
		assert.Equal(t, input.Username, res.Username, "Nama user harus sesuai")
		repo.AssertExpectations(t)
	})
	t.Run("Duplicated", func(t *testing.T) {
		repo.On("Insert", mock.Anything).Return(domain.Core{}, errors.New(config.DUPLICATED_DATA)).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Email: "skjsa", Password: "same"}
		res, err := srv.AddUser(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, config.DUPLICATED_DATA, "Pesan error tidak sesuai")
		assert.Equal(t, uint(0), res.ID, "ID seharusnya 0 karena gagal input data")
		repo.AssertExpectations(t)
	})
	t.Run("Gaga; Hash", func(t *testing.T) {
		repo.On("Insert", mock.Anything).Return(domain.Core{}, errors.New("cannot encript password")).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Email: "skjsa", Password: ""}
		res, err := srv.AddUser(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "cannot encript password", "Pesan error tidak sesuai")
		assert.Equal(t, uint(0), res.ID, "ID seharusnya 0 karena gagal input data")
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Login", func(t *testing.T) {
		repo.On("Login", mock.Anything).Return(domain.Core{ID: uint(1), Username: "same"}, nil).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Password: "same"}
		res, err := srv.LoginUser(input)
		assert.Nil(t, err)
		assert.NotEmpty(t, res.ID, "Seharusnya ada ID yang berhasil dibuat")
		assert.NotEqual(t, res.Password, input.Password, "Password tidak terenkripsi")
		assert.Equal(t, input.Username, res.Username, "Nama user harus sesuai")
		repo.AssertExpectations(t)
	})
	t.Run("Gagal Login", func(t *testing.T) {
		repo.On("Login", mock.Anything).Return(domain.Core{}, errors.New("cannot encript password")).Once()
		srv := New(repo)
		input := domain.Core{Username: "same", Password: ""}
		res, err := srv.LoginUser(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "cannot encript password", "Pesan error tidak sesuai")
		assert.Equal(t, uint(0), res.ID, "ID seharusnya 0 karena gagal input data")
		repo.AssertExpectations(t)
	})
}
