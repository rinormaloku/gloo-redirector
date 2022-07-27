package mesh

import (
	"bytes"
	_ "embed"
	"github.com/rinormaloku/gloo-redirector/pkg/domain"
	"github.com/rinormaloku/gloo-redirector/pkg/templates"
)

func GenerateEdgeRedirections(inputData domain.InputData) (bytes.Buffer, error) {
	inputData.DefaultTemplate = templates.GlooEdgeResourceTemplate
	return generate(inputData)
}
