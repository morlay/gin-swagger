package codegen

func NewPrinter() *Printer {
	return &Printer{
		Indent:    "        ",
		SpaceChar: " ",
		LineBreak: "\n",
	}
}

type Printer struct {
	SpaceChar string
	Indent    string
	LineBreak string
	code      string
}

func (p *Printer) Input(s string) *Printer {
	p.code = p.code + s
	return p
}

func (p *Printer) Space() *Printer {
	return p.Input(p.SpaceChar)
}

func (p *Printer) Tab() *Printer {
	return p.Input(p.Indent)
}

func (p *Printer) NewLine() *Printer {
	return p.Input(p.LineBreak)
}

func (p *Printer) String() string {
	return p.code + ""
}
