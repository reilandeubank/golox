package parser

type ExprVisitor interface {
	VisitBinaryExpr(b Binary) (interface{}, error)
	VisitGroupingExpr(g Grouping) (interface{}, error)
	VisitLiteralExpr(l Literal) (interface{}, error)
	VisitUnaryExpr(u Unary) (interface{}, error)
	//VisitVariableExpr(v *Variable) interface{}
	//VisitAssignExpr(a *Assign) interface{}
	//VisitLogicalExpr(l *Logical) interface{}
	//VisitCallExpr(c *Call) interface{}
	// VisitGetExpr(g *Get) interface{}
	// VisitSetExpr(s *Set) interface{}
	// VisitThisExpr(t *This) interface{}
	// VisitSuperExpr(s *Super) interface{}
}