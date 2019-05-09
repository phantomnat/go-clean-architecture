// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import article "github.com/phantomnat/go-clean-architecture/article"
import mock "github.com/stretchr/testify/mock"

// ArticleRepository is an autogenerated mock type for the ArticleRepository type
type ArticleRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *ArticleRepository) Delete(id int64) (bool, error) {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Fetch provides a mock function with given fields: cursor, num
func (_m *ArticleRepository) Fetch(cursor string, num int64) ([]*article.Article, error) {
	ret := _m.Called(cursor, num)

	var r0 []*article.Article
	if rf, ok := ret.Get(0).(func(string, int64) []*article.Article); ok {
		r0 = rf(cursor, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*article.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(cursor, num)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *ArticleRepository) GetByID(id int64) (*article.Article, error) {
	ret := _m.Called(id)

	var r0 *article.Article
	if rf, ok := ret.Get(0).(func(int64) *article.Article); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*article.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: title
func (_m *ArticleRepository) GetByTitle(title string) (*article.Article, error) {
	ret := _m.Called(title)

	var r0 *article.Article
	if rf, ok := ret.Get(0).(func(string) *article.Article); ok {
		r0 = rf(title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*article.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ar
func (_m *ArticleRepository) Store(ar *article.Article) (int64, error) {
	ret := _m.Called(ar)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*article.Article) int64); ok {
		r0 = rf(ar)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*article.Article) error); ok {
		r1 = rf(ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ar
func (_m *ArticleRepository) Update(ar *article.Article) (*article.Article, error) {
	ret := _m.Called(ar)

	var r0 *article.Article
	if rf, ok := ret.Get(0).(func(*article.Article) *article.Article); ok {
		r0 = rf(ar)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*article.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*article.Article) error); ok {
		r1 = rf(ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}