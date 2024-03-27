package recognizer

import "context"

type Service interface {
	RecognizeFormula(img []byte, ctx context.Context) (string, error)
	RecognizeText(imgs [][]byte, ctx context.Context) (string, error)
	RecognizeMixed(imgs [][]byte, ctx context.Context) (string, error)
}
