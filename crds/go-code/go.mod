module github.com/lyraproj/crd-mod

require (
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v0.0.0-20160705203006-01aeca54ebda // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gophercloud/gophercloud v0.0.0-20190126172459-c818fa66e4c8 // indirect
	github.com/lyraproj/lyra-operator v0.0.0-20190412150939-82bb153789bc
	github.com/operator-framework/operator-sdk v0.7.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	google.golang.org/appengine v1.5.0 // indirect
	k8s.io/api v0.0.0-20190511023547-e63b5755afac
	k8s.io/apiextensions-apiserver v0.0.0-20190408173516-34b98ea6c731
	k8s.io/apimachinery v0.0.0-20190511023455-ad85901afca0

	k8s.io/client-go v10.0.0+incompatible
	k8s.io/sample-controller v0.0.0-20190511024443-0d5581c413c7
	sigs.k8s.io/controller-runtime v0.1.10
)

replace k8s.io/api => k8s.io/api v0.0.0-20181221193117-173ce66c1e39

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190126155707-0e6dcdd1b5ce

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190103235604-e7617803aceb
