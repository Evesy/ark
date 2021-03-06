package restic

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"k8s.io/client-go/kubernetes"

	"github.com/heptio/ark/pkg/client"
	clientset "github.com/heptio/ark/pkg/generated/clientset/versioned"
	arktest "github.com/heptio/ark/pkg/util/test"
)

type fakeFactory struct{}

var _ client.Factory = &fakeFactory{}

func (f *fakeFactory) BindFlags(flags *pflag.FlagSet) {
	panic("not implemented")
}

func (f *fakeFactory) Client() (clientset.Interface, error) {
	panic("not implemented")
}

func (f *fakeFactory) KubeClient() (kubernetes.Interface, error) {
	panic("not implemented")
}

func (f *fakeFactory) Namespace() string {
	return ""
}

func TestComplete(t *testing.T) {
	// no key options provided should error
	o := &InitRepositoryOptions{}
	err := o.Complete(&fakeFactory{})
	assert.EqualError(t, err, errKeySizeTooSmall.Error())

	// both KeyFile and KeyData provided should error
	o = &InitRepositoryOptions{
		KeyFile: "/foo",
		KeyData: "bar",
	}
	err = o.Complete(&fakeFactory{})
	assert.EqualError(t, err, errKeyFileAndKeyDataProvided.Error())

	// if KeyFile is provided, its contents are used
	fileContents := []byte("bar")
	o = &InitRepositoryOptions{
		KeyFile:    "/foo",
		fileSystem: arktest.NewFakeFileSystem().WithFile("/foo", fileContents),
	}
	assert.NoError(t, o.Complete(&fakeFactory{}))
	assert.Equal(t, fileContents, o.keyBytes)

	// if KeyData is provided, it's used
	o = &InitRepositoryOptions{
		KeyData: "bar",
	}
	assert.NoError(t, o.Complete(&fakeFactory{}))
	assert.Equal(t, []byte(o.KeyData), o.keyBytes)

	// if KeySize is provided, a random key is generated
	o = &InitRepositoryOptions{
		KeySize: 10,
	}
	assert.NoError(t, o.Complete(&fakeFactory{}))
	assert.Equal(t, o.KeySize, len(o.keyBytes))
}
