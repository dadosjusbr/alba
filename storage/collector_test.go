package storage

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestConnect_EnvVarNotDefined(t *testing.T) {
	is := is.New(t)
	os.Unsetenv("MONGODB")
	client, err := conect()
	is.True(client == nil)
	is.True(err.Error() == "error trying get environment variable:\"$MONGODB is empty\"")
}
