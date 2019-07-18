package apis

import (
	accessv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/access/v1alpha1"
	specsv1alpha1 "github.com/deislabs/smi-sdk-go/pkg/apis/specs/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, accessv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, specsv1alpha1.SchemeBuilder.AddToScheme)
}
