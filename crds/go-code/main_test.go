package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	testSingular = "widget"
	testPlural   = "widgets"
	testGroup    = "lyra.example.com"
)

func TestCreateCRD(t *testing.T) {
	var err error
	var crds []string

	//delete just in case
	err = deleteCRD(testGroup, testSingular, testPlural)
	require.NoError(t, err)

	//confirm not present
	crds, err = listCRDs()
	require.NoError(t, err)
	require.NotContains(t, crds, crdName(testPlural, testGroup))

	//create
	err = createCRD(testGroup, testSingular, testPlural)
	require.NoError(t, err)

	//check present
	crds, err = listCRDs()
	require.NoError(t, err)
	require.Contains(t, crds, crdName(testPlural, testGroup))

	//delete
	err = deleteCRD(testGroup, testSingular, testPlural)
	require.NoError(t, err)

	//have a little hacky sleep
	time.Sleep(1 * time.Second)

	//check absent again
	crds, err = listCRDs()
	require.NoError(t, err)
	require.NotContains(t, crds, crdName(testPlural, testGroup))
}
