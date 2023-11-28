package parser

type ExprVisitor interface {
	VisitBinaryExpr(b Binary) (interface{}, error)
	VisitGroupingExpr(g Grouping) (interface{}, error)
	VisitLiteralExpr(l Literal) (interface{}, error)
	VisitUnaryExpr(u Unary) (interface{}, error)
	VisitVariableExpr(v Variable) (interface{}, error)
	VisitAssignExpr(a Assign) (interface{}, error)
	//VisitLogicalExpr(l Logical) (interface{}, error)
	//VisitCallExpr(c Call) (interface{}, error)
	// VisitGetExpr(g Get) (interface{}, error)
	// VisitSetExpr(s Set) (interface{}, error)
	// VisitThisExpr(t This) (interface{}, error)
	// VisitSuperExpr(s Super) (interface{}, error)
}

type StmtVisitor interface {
	VisitExprStmt(e ExprStmt) (interface{}, error)
	VisitPrintStmt(p PrintStmt) (interface{}, error)
	VisitVarStmt(v VarStmt) (interface{}, error)
	VisitBlockStmt(b BlockStmt) (interface{}, error)
	// VisitIfStmt(i IfStmt) (interface{}, error)
	// VisitWhileStmt(w WhileStmt) (interface{}, error)
	// VisitFunStmt(f FunStmt) (interface{}, error)
	// VisitReturnStmt(r ReturnStmt) (interface{}, error)
	// VisitClassStmt(c ClassStmt) (interface{}, error)
}