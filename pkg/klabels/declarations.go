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
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type FieldDeclaration struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Annotations string `json:"annotation"`
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
	FunctionKind         LabelKind = "Function"
	VariableKind                   = "Variable"
	TypeDeclarationKind            = "Type"
	FieldDeclarationKind           = "Field"
	ConstDeclarationKind           = "Const"
	VarDeclarationKind             = "Var"
	UnknownKind                    = "Unknown"
	CommentKind                    = "Comment"
	PackageKind                    = "Package"
	ReturnKind                     = "ReturnStatement"
	ScopeBeginKind                 = "ScopeBegin"
	ScopeEndKind                   = "ScopeEnd"
	ConstScopeBeginKind            = "ConstantScopeBegin"
	VarScopeBeginKind              = "VarScopeBegin"
	ConstScopeEndKind              = "ConstantScopeEnd"
	VarScopeEndKind                = "VarScopeEnd"
)
