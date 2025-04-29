package assets

import (
	"embed"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	//go:embed manifests/data_nodes_stateful_set.yaml
	dataNodeStatefulSet embed.FS

	//go:embed manifests/name_nodes_deployment.yaml
	nameNodeDeployment embed.FS

	appsScheme = runtime.NewScheme()

	appsCodecs = serializer.NewCodecFactory(appsScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
}

func GetDataNodeStatefulSetFromFile() *appsv1.StatefulSet {
	bytes, err := nameNodeDeployment.ReadFile("manifests/data_nodes_stateful_set.yaml")
	if err != nil {
		panic(err)
	}

	o, err := runtime.Decode(appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion), bytes)
	if err != nil {
		panic(err)
	}

	s, ok := o.(*appsv1.StatefulSet)
	if !ok {
		panic("decoded object is not a StatefulSet")
	}

	return s
}

func GetNameNodeDeploymentFromFile() *appsv1.Deployment {
	bytes, err := dataNodeStatefulSet.ReadFile("manifests/name_nodes_deployment.yaml")
	if err != nil {
		panic(err)
	}

	o, err := runtime.Decode(appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion), bytes)
	if err != nil {
		panic(err)
	}

	d, ok := o.(*appsv1.Deployment)
	if !ok {
		panic("decoded object is not a Deployment")
	}

	return d
}
