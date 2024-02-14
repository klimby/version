package cmd

import "github.com/stretchr/testify/mock"

type __runnerMock struct {
	mock.Mock
}

func (r *__runnerMock) Run() error {
	ret := r.Called()

	return ret.Error(0)
}

func __newRunnerMock(err error) *__runnerMock {
	m := &__runnerMock{}
	m.On("Run").Return(err)

	return m
}

type __helpMock struct {
	mock.Mock
}

func (h *__helpMock) Help() {
	h.Called()
}

func __newHelpMock() *__helpMock {
	m := &__helpMock{}
	m.On("Help")

	return m
}
