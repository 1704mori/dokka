package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocked struct {
	mock.Mock
	DockerAPI
}

func (m *mocked) ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called()
	containers, ok := args.Get(0).([]types.Container)
	if !ok && args.Get(0) != nil {
		panic("Containers is not of type []types.Container")
	}

	return containers, args.Error(1)
}

func Test_Docker_ListContainers_nil(t *testing.T) {
	m := new(mocked)
	m.On("ContainerList", mock.Anything, mock.Anything).Return(nil, nil)
	cli := &Client{
		cli: m,
	}

	containers, err := cli.ListContainers()
	assert.Empty(t, containers, "containers should be empty")
	require.NoError(t, err, "should not return any error")

	m.AssertExpectations(t)
}

func Test_Docker_ListContainers_error(t *testing.T) {
	m := new(mocked)
	m.On("ContainerList", mock.Anything, mock.Anything).Return(nil, errors.New("test"))
	cli := &Client{
		cli: m,
	}

	containers, err := cli.ListContainers()
	assert.Empty(t, containers, "containers should be empty")
	require.Error(t, err, "should return 'test' error")

	m.AssertExpectations(t)
}

func Test_Docker_ListContainers_success(t *testing.T) {
	mockedContainers := []types.Container{
		{
			ID:    "dso8fhnlask9_ayo",
			Names: []string{"i_am_a_container"},
		},
		{
			ID:    "4d8fg4ds5fg4wasd_tokyo",
			Names: []string{"a_lonely_container"},
		},
	}

	m := new(mocked)
	m.On("ContainerList", mock.Anything, mock.Anything).Return(mockedContainers, nil)
	cli := &Client{
		cli: m,
	}

	containers, err := cli.ListContainers()

	require.NoError(t, err, "should not return any error")
	assert.Equal(t, containers, mockedContainers)

	m.AssertExpectations(t)
}

func Test_Docker_FindContainer_error(t *testing.T) {
	mockedContainers := []types.Container{
		{
			ID:    "dso8fhnlask9_ayo",
			Names: []string{"i_am_a_container"},
		},
		{
			ID:    "4d8fg4ds5fg4wasd_tokyo",
			Names: []string{"a_lonely_container"},
		},
	}

	m := new(mocked)
	m.On("ContainerList", mock.Anything, mock.Anything).Return(mockedContainers, nil)

	cli := &Client{
		cli: m,
	}

	_, err := cli.FindContainer("invalid_id")
	require.Error(t, err, "should return an error")

	m.AssertExpectations(t)
}

func Test_Docker_FindContainer_success(t *testing.T) {
	mockedContainers := []types.Container{
		{
			ID:    "dso8fhnlask9_ayo",
			Names: []string{"i_am_a_container"},
		},
		{
			ID:    "4d8fg4ds5fg4wasd_tokyo",
			Names: []string{"a_lonely_container"},
		},
	}

	m := new(mocked)
	m.On("ContainerList", mock.Anything, mock.Anything).Return(mockedContainers, nil)

	cli := &Client{
		cli: m,
	}

	container, err := cli.FindContainer("dso8fhnlask9_ayo")
	require.NoError(t, err, "should not return an error")

	assert.Equal(t, container, mockedContainers[0])

	m.AssertExpectations(t)
}
