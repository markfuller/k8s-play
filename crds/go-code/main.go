package main

import (
	"fmt"
	"strings"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	fmt.Printf("hello world\n")
	fmt.Printf("hello world %v", v1alpha1.Kind("foo"))
}

func getClient() (apiextensionsv1beta1.ApiextensionsV1beta1Interface, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return apiextensionsv1beta1.NewForConfig(cfg)
}

func listCRDs() ([]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	list, err := client.CustomResourceDefinitions().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	crdNames := []string{}
	for _, v := range list.Items {
		crdNames = append(crdNames, v.Name)
	}
	fmt.Printf("Returning %v\n", crdNames)
	return crdNames, nil
}

func deleteCRD(group string, singular string, plural string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	name := crdName(plural, group)
	err = client.CustomResourceDefinitions().Delete(name, &metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func createCRD(group string, singular string, plural string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	name := crdName(plural, group)
	crd := &apiext.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: apiext.CustomResourceDefinitionSpec{
			Group:   group,
			Version: "v1beta1",
			Names: apiext.CustomResourceDefinitionNames{
				Plural:   plural,
				Singular: singular,
				Kind:     strings.Title(singular),
			},
		},
	}
	_, err = client.CustomResourceDefinitions().Create(crd)

	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func crdName(plural string, group string) string {
	return fmt.Sprintf("%s.%s", plural, group)
}
