package calculo_imposto

import "github.com/chonla/format"

type ExplainIntf interface {
	Add(source string, periodo int, tipoPeriodo TipoPeriodo, message string, params map[string]interface{})
	Items() []*ExplainItem
}

type Explain struct {
	items []*ExplainItem
}

func NewExplain() *Explain {
	return &Explain{}
}

func (e *Explain) Add(source string, periodo int, tipoPeriodo TipoPeriodo, message string, params map[string]interface{}) {
	e.items = append(e.items, &ExplainItem{
		Source:      source,
		Periodo:     periodo,
		TipoPeriodo: tipoPeriodo,
		Message:     message,
		Params:      params,
	})
}

func (e *Explain) Items() []*ExplainItem {
	return e.items
}

type ExplainItem struct {
	Source      string                 `json:"source"`
	Periodo     int                    `json:"periodo"`
	TipoPeriodo TipoPeriodo            `json:"tipo-periodo"`
	Message     string                 `json:"message"`
	Params      map[string]interface{} `json:"params"`
}

func (e *ExplainItem) FormatMessage() string {
	return format.Sprintf(e.Message, e.Params)
}

type ExplainEmpty struct{}

func (e *ExplainEmpty) Add(string, int, TipoPeriodo, string, map[string]interface{}) {

}

func (e *ExplainEmpty) Items() []*ExplainItem {
	return nil
}
