package templates

import _ "embed"

//go:embed mesh-resources.yaml
var GlooMeshResourceTemplate string

//go:embed edge-resources.yaml
var GlooEdgeResourceTemplate string
