package docgen

import "errors"

var (
	ErrTemplateNotFound = errors.New("docgen: template not found")
	ErrRendererNotFound = errors.New("docgen: renderer not found")
	ErrInvalidTemplate  = errors.New("docgen: invalid template")
	ErrRenderFailed     = errors.New("docgen: render failed")
)
