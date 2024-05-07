package domain

const (
	WeightRequestWordIndex     = 0.3
	WeightRequestWordDuplicate = 1.0

	WeightComicsWordIndex              = 0.01
	WeightComicsWordDuplicate          = 0.5
	WeightComicsActual                 = 0.001
	WeightComicsWordPositionTitle      = 0.6
	WeightComicsWordPositionTranscript = 0.1
	WeightComicsWordPositionAlt        = 0.2

	WeightResponseRelevantCoverage           = 5
	WeightResponseRelevantCoverageTitle      = 1
	WeightResponseRelevantCoverageTranscript = 1
	WeightResponseRelevantCoverageAlt        = 1
)

type Weights struct {
	Word     *Words
	Comic    *Comics
	Position *Positions
	Weight   float64
}
