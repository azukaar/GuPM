package main

import (
	"testing"
	// "reflect"
	_ "fmt"
	"github.com/stretchr/testify/require"
)

func TestArgs(t *testing.T) {
	require := require.New(t)

	// Basic
	command, args := GetArgs([]string{"make", "-p", "npm"})
	require.Equal("make", command)
	require.Equal(Arguments{"$0": "make --p npm", "P": "npm"}, args)

	// Nothing
	command, args = GetArgs([]string{})
	require.Equal("", command)
	require.Equal(Arguments{"$0": ""}, args)

	// No args
	command, args = GetArgs([]string{"install"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install"}, args)

	// bools and long mix
	command, args = GetArgs([]string{"install", "--dev", "--provider", "npm"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install --dev --provider npm", "Dev": "true", "Provider": "npm"}, args)
	command, args = GetArgs([]string{"install", "--provider", "npm", "--dev"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install --provider npm --dev", "Dev": "true", "Provider": "npm"}, args)

	// bools at the end
	command, args = GetArgs([]string{"install", "--dev"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install --dev", "Dev": "true"}, args)

	// long args
	command, args = GetArgs([]string{"install", "--provider", "npm"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install --provider npm", "Provider": "npm"}, args)

	// anonymous
	command, args = GetArgs([]string{"install", "npm"})
	require.Equal("install", command)
	require.Equal(Arguments{"$0": "install npm", "$1": "npm"}, args)

	// asList
	command, args = GetArgs([]string{"install", "go", "npm"})
	require.Equal([]string{"go", "npm"}, args.AsList())

}
