package klabels

type Analysis struct {
	Labels []Label `json:"labels"`
	Line   int     `json:"line"`
}

type Label struct {
	Kind LabelKind `json:"kind"`
	Data any       `json:"data"`
	Line int       `json:"line"`
}

type FunctionDeclaration struct {
	Name       string              `json:"name"`
	Parameters []FunctionParameter `json:"parameters"`
	Result     []string            `json:"result"`
}

type FunctionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type VariableDeclaration struct {
	Name         string  `json:"name"`
	Type         *string `json:"type"`
	Value        *string `json:"value"`
	Reassignment bool    `json:"reassignment"`
}

type TypeDeclaration struct {
	Name       string         `json:"name"`
	Kind       string         `json:"kind"`
	Properties []TypeProperty `json:"properties"`
}

type TypeProperty struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Annotation string `json:"annotation"`
}

type ConstantDeclaration struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value *string `json:"value"`
}

type PackageDeclaration struct {
	Name string
}

type ReturnStatement struct {
	Values []string
}

type LabelKind string
type TypeKind int

const (
	FunctionKind        LabelKind = "Function"
	VariableKind                  = "Variable"
	ConstantKind                  = "Constant"
	TypeDeclarationKind           = "Type"
	UnknownKind                   = "Unknown"
	CommentKind                   = "Comment"
	PackageKind                   = "Package"
	ReturnKind                    = "ReturnStatement"
	ScopeBeginKind                = "ScopeBegin"
	ScopeEndKind                  = "ScopeEnd"
)
