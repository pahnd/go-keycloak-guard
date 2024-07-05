package permission

import (
	"errors"
)

const StrategyAffirmative = "affirmative"
const StrategyConsensus = "consensus"
const StrategyUnanimous = "unanimous"

type DecisionStrategyClient struct {
	affirmative *AffirmativeDecisionAStrategy
	consensus   *ConsensusDecisionStrategy
	unanimous   *UnanimousDecisionStrategy
}

func NewDecisionStrategyClient() *DecisionStrategyClient {
	return &DecisionStrategyClient{
		affirmative: &AffirmativeDecisionAStrategy{},
		consensus:   &ConsensusDecisionStrategy{},
		unanimous:   &UnanimousDecisionStrategy{},
	}
}

func (d *DecisionStrategyClient) getDecision(strategy string) (DecisionInterface, error) {
	switch strategy {
	case StrategyAffirmative:
		return d.affirmative, nil
	case StrategyConsensus:
		return d.consensus, nil
	case StrategyUnanimous:
		return d.unanimous, nil
	default:
		return nil, errors.New("invalid strategy name")
	}
}

func (d *DecisionStrategyClient) HasPermissions(requiredPermissions []string, permissions PermissionCollection, strategy string) (bool, error) {
	decision, err := d.getDecision(strategy)
	if err != nil {
		return false, err
	}
	return decision.HasPermissions(requiredPermissions, permissions), nil
}
