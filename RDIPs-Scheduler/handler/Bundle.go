package handler

import (
	"fmt"

	api "github.com/evanw/esbuild/pkg/api"
	"github.com/t14raptor/go-fast/ast"
	"github.com/t14raptor/go-fast/parser"
)

type FunctionMetrics struct {
	cyclomaticComplexity int
	loopComplexity       float64
	recursionDepth       int
	callCount            int64
	nestingLevel         int
}

type TraversalContext struct {
	functionName  string
	inLoop        bool
	inConditional bool
	nestingLevel  int
}

type ComplexityAnalyzer struct {
	functions map[string]*FunctionMetrics
}

func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		functions: make(map[string]*FunctionMetrics),
	}
}

func (ca *ComplexityAnalyzer) ensureFunction(name string) *FunctionMetrics {
	if name == "" {
		name = "*__default"
	}

	if f, ok := ca.functions[name]; ok {
		return f
	}

	if name == "*__default" {
		ca.functions[name] = &FunctionMetrics{
			cyclomaticComplexity: 1, // Base complexity
			loopComplexity:       0,
			recursionDepth:       1, // Start with 1 for the function itself
			callCount:            1,
			nestingLevel:         0,
		}
	} else {
		ca.functions[name] = &FunctionMetrics{
			cyclomaticComplexity: 1, // Base complexity
			loopComplexity:       0,
			recursionDepth:       1, // Start with 1 for the function itself
			callCount:            0,
			nestingLevel:         0,
		}
	}
	return ca.functions[name]
}

func (ca *ComplexityAnalyzer) traverse(node ast.Node, context *TraversalContext) {
	if node == nil {
		return
	}

	// Ensure we have a function context
	currentFunction := ca.ensureFunction(context.functionName)

	switch n := node.(type) {
	case *ast.FunctionDeclaration:
		if ca.safeFunctionDeclaration(n) {
			functionName := n.Function.Name.Name
			newContext := &TraversalContext{
				functionName:  functionName,
				inLoop:        false,
				inConditional: false,
				nestingLevel:  0,
			}
			ca.ensureFunction(functionName)
			ca.traverse(n.Function.Body, newContext)
		}

	case *ast.ReturnStatement:
		if n.Argument != nil {
			ca.traverse(n.Argument, context)
		}

	case *ast.BlockStatement:
		for _, stmt := range n.List {
			ca.traverse(stmt.Stmt, context)
		}

	case *ast.CatchStatement:
		currentFunction.cyclomaticComplexity++
		if n.Body != nil {
			ca.traverse(n.Body, context)
		}

	case *ast.TryStatement:
		if n.Body != nil {
			ca.traverse(n.Body, context)
		}
		if n.Catch != nil {
			ca.traverse(n.Catch, context)
		}
		if n.Finally != nil {
			ca.traverse(n.Finally, context)
		}

	case *ast.CallExpression:
		ca.handleCallExpression(n, context)

	case *ast.IfStatement:
		ca.handleIfStatement(n, context)

	case *ast.ForStatement:
		ca.handleLoopStatement(n.Body, n.Test, context, "for")

	case *ast.WhileStatement:
		ca.handleLoopStatement(n.Body, n.Test, context, "while")

	case *ast.DoWhileStatement:
		ca.handleLoopStatement(n.Body, n.Test, context, "dowhile")

	case *ast.SwitchStatement:
		ca.handleSwitchStatement(n, context)

	case *ast.ConditionalExpression:
		ca.handleConditionalExpression(n, context)

	case *ast.BinaryExpression:
		ca.traverse(n.Left, context)
		ca.traverse(n.Right, context)

	case *ast.AssignExpression:
		ca.traverse(n.Left, context)
		ca.traverse(n.Right, context)

	case *ast.ExpressionStatement:
		ca.traverse(n.Expression, context)

	case *ast.MemberExpression:
		ca.traverse(n.Object, context)
	case *ast.Expression:
		if n.Expr != nil {
			ca.traverse(n.Expr, context)
		}
	case *ast.SequenceExpression:
		for _, expr := range n.Sequence {
			ca.traverse(&expr, context)
		}
	case *ast.Statement:
		ca.traverse(n.Stmt, context)
	default:
	}
}

func (ca *ComplexityAnalyzer) safeFunctionDeclaration(n *ast.FunctionDeclaration) bool {
	return n != nil && n.Function != nil && n.Function.Name != nil && n.Function.Body != nil
}

func (ca *ComplexityAnalyzer) handleCallExpression(n *ast.CallExpression, context *TraversalContext) {
	if n.Callee != nil && n.Callee.Expr != nil {
		if ident, ok := n.Callee.Expr.(*ast.Identifier); ok && ident != nil {
			currentFunction := ca.ensureFunction(ident.Name)

			// Check for recursion
			if ident.Name == context.functionName {
				currentFunction.recursionDepth += 2
				// Add complexity for recursion (linear increase, not exponential)
				currentFunction.cyclomaticComplexity++
			} else {
				currentFunction.callCount++
			}
		}

		if member, ok := n.Callee.Expr.(*ast.MemberExpression); ok && member != nil {
			ca.traverse(member, context)
		}

		args := n.ArgumentList
		for _, arg := range args {
			ca.traverse(arg.Expr, context)
		}
	}
}

func (ca *ComplexityAnalyzer) handleIfStatement(n *ast.IfStatement, context *TraversalContext) {
	currentFunction := ca.ensureFunction(context.functionName)
	currentFunction.cyclomaticComplexity++

	// Apply nesting penalty correctly
	if context.inLoop {
		currentFunction.loopComplexity += 0.5 // Penalty for conditional inside loop
	}
	if context.inConditional {
		currentFunction.cyclomaticComplexity++ // Penalty for nested conditional
	}

	// Create new context for the conditional body
	newContext := &TraversalContext{
		functionName:  context.functionName,
		inLoop:        context.inLoop,
		inConditional: true,
		nestingLevel:  context.nestingLevel + 1,
	}

	ca.traverse(n.Consequent, newContext)
	if n.Alternate != nil {
		ca.traverse(n.Alternate, newContext)
	}
}

func (ca *ComplexityAnalyzer) handleLoopStatement(body ast.Node, test ast.Node, context *TraversalContext, loopType string) {
	currentFunction := ca.ensureFunction(context.functionName)
	currentFunction.cyclomaticComplexity++

	// Calculate loop complexity based on nesting
	loopComplexity := 1.0
	if context.inLoop {
		loopComplexity = 2.0 // Nested loop penalty
	}
	if context.inConditional {
		loopComplexity += 0.5 // Loop inside conditional penalty
	}

	currentFunction.loopComplexity += loopComplexity

	// Create new context for loop body
	newContext := &TraversalContext{
		functionName:  context.functionName,
		inLoop:        true,
		inConditional: context.inConditional,
		nestingLevel:  context.nestingLevel + 1,
	}
	if test != nil {
		ca.traverse(test, newContext)
	}
	if body != nil {
		ca.traverse(body, newContext)
	}
}

func (ca *ComplexityAnalyzer) handleSwitchStatement(n *ast.SwitchStatement, context *TraversalContext) {
	currentFunction := ca.ensureFunction(context.functionName)

	// Switch adds base complexity + complexity for each case
	if n.Body != nil {
		currentFunction.cyclomaticComplexity += len(n.Body) - 1 // Each case adds complexity

		// Apply nesting penalty
		if context.inLoop {
			currentFunction.loopComplexity += 0.5
		}

		newContext := &TraversalContext{
			functionName:  context.functionName,
			inLoop:        context.inLoop,
			inConditional: true,
			nestingLevel:  context.nestingLevel + 1,
		}

		for _, c := range n.Body {
			for _, stmt := range c.Consequent {
				ca.traverse(stmt.Stmt, newContext)
			}
		}
	}
}

func (ca *ComplexityAnalyzer) handleConditionalExpression(n *ast.ConditionalExpression, context *TraversalContext) {
	currentFunction := ca.ensureFunction(context.functionName)
	currentFunction.cyclomaticComplexity++

	// Apply nesting penalty
	if context.inLoop {
		currentFunction.loopComplexity += 0.5
	}
	if context.inConditional {
		currentFunction.cyclomaticComplexity++
	}

	newContext := &TraversalContext{
		functionName:  context.functionName,
		inLoop:        context.inLoop,
		inConditional: true,
		nestingLevel:  context.nestingLevel + 1,
	}

	ca.traverse(n.Test, newContext)
	ca.traverse(n.Consequent, newContext)
	ca.traverse(n.Alternate, newContext)
}

func (ca *ComplexityAnalyzer) analyzeProgram(program *ast.Program) error {
	if program == nil || program.Body == nil {
		return fmt.Errorf("invalid program structure")
	}

	// Create initial context for top-level code
	initialContext := &TraversalContext{
		functionName:  "*__default",
		inLoop:        false,
		inConditional: false,
		nestingLevel:  0,
	}

	for _, stmt := range program.Body {
		ca.traverse(stmt.Stmt, initialContext)
	}

	return nil
}

func (ca *ComplexityAnalyzer) getMetrics() map[string]float64 {
	metricsMap := map[string]float64{
		"totalFunctionComplexity":     0,
		"totalConditionalComplexity":  0,
		"totalLoopComplexity":         0,
		"totalCallCount":              0,
		"totalRecursionDepth":         0,
		"averageCyclomaticComplexity": 0,
	}

	for _, metrics := range ca.functions {
		funcComplexity := (float64(metrics.callCount)*float64(metrics.recursionDepth))*float64(metrics.cyclomaticComplexity) + metrics.loopComplexity
		metricsMap["totalFunctionComplexity"] += funcComplexity

		metricsMap["totalConditionalComplexity"] += float64(metrics.cyclomaticComplexity)
		metricsMap["totalLoopComplexity"] += metrics.loopComplexity
		metricsMap["totalCallCount"] += float64(metrics.callCount)
		metricsMap["totalRecursionDepth"] += float64(metrics.recursionDepth)
		metricsMap["averageCyclomaticComplexity"] += float64(metrics.cyclomaticComplexity)

	}

	return metricsMap
}

func BundleJSCode(jsCode string) ([]byte, error) {
	result := api.Transform(jsCode, api.TransformOptions{
		Loader:            api.LoaderJS,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
	})
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("%s", result.Errors[0].Text)
	}
	return result.Code, nil
}

func AstParser(code string) (float64, map[string]float64, error) {
	if code == "" {
		return 0, nil, fmt.Errorf("empty code provided")
	}

	program, err := parser.ParseFile(code)
	if err != nil {
		return 0, nil, fmt.Errorf("parsing error: %w", err)
	}

	analyzer := NewComplexityAnalyzer()
	if err := analyzer.analyzeProgram(program); err != nil {
		return 0, nil, fmt.Errorf("analysis error: %w", err)
	}

	metricsMap := analyzer.getMetrics()
	return float64(len(code)) / 1024, metricsMap, nil
}
