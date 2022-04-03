package agents

type storage interface {
	UpdateAgentConn(id string, conn string) error
	GetAgent(id string) (*Agent, error)
}

type AgentController struct {
	db storage
}

func New(storage storage) *AgentController {
	return &AgentController{db: storage}
}

func (am *AgentController) Register(agent *Agent) error {
	return nil
}

func (am *AgentController) Report(aID string, conn string) error {
	return am.db.UpdateAgentConn(aID, conn)
}

func (am *AgentController) GetAgent(id string) (*Agent, error) {
	return am.db.GetAgent(id)
}
