package mesh

import (
	"bytes"
	_ "embed"
	"github.com/rinormaloku/gloo-redirector/pkg/domain"
	"github.com/rinormaloku/gloo-redirector/pkg/templates"
)

func GenerateMeshRedirections(inputData domain.InputData) (bytes.Buffer, error) {
	inputData.DefaultTemplate = templates.GlooMeshResourceTemplate
	return generate(inputData)
}
