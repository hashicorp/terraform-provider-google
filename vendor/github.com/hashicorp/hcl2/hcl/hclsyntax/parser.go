package hclsyntax

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/apparentlymart/go-textseg/textseg"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type parser struct {
	*peeker

	// set to true if any recovery is attempted. The parser can use this
	// to attempt to reduce error noise by suppressing "bad token" errors
	// in recovery mode, assuming that the recovery heuristics have failed
	// in this case and left the peeker in a wrong place.
	recovery bool
}

func (p *parser) ParseBody(end TokenType) (*Body, hcl.Diagnostics) {
	attrs := Attributes{}
	blocks := Blocks{}
	var diags hcl.Diagnostics

	startRange := p.PrevRange()
	var endRange hcl.Range

Token:
	for {
		next := p.Peek()
		if next.Type == end {
			endRange = p.NextRange()
			p.Read()
			break Token
		}

		switch next.Type {
		case TokenNewline:
			p.Read()
			continue
		case TokenIdent:
			item, itemDiags := p.ParseBodyItem()
			diags = append(diags, itemDiags...)
			switch titem := item.(type) {
			case *Block:
				blocks = append(blocks, titem)
			case *Attribute:
				if existing, exists := attrs[titem.Name]; exists {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute redefined",
						Detail: fmt.Sprintf(
							"The attribute %q was already defined at %s. Each attribute may be defined only once.",
							titem.Name, existing.NameRange.String(),
						),
						Subject: &titem.NameRange,
					})
				} else {
					attrs[titem.Name] = titem
				}
			default:
				// This should never happen for valid input, but may if a
				// syntax error was detected in ParseBodyItem that prevented
				// it from even producing a partially-broken item. In that
				// case, it would've left at least one error in the diagnostics
				// slice we already dealt with above.
				//
				// We'll assume ParseBodyItem attempted recovery to leave
				// us in a reasonable position to try parsing the next item.
				continue
			}
		default:
			bad := p.Read()
			if !p.recovery {
				if bad.Type == TokenOQuote {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid attribute name",
						Detail:   "Attribute names must not be quoted.",
						Subject:  &bad.Range,
					})
				} else {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute or block definition required",
						Detail:   "An attribute or block definition is required here.",
						Subject:  &bad.Range,
					})
				}
			}
			endRange = p.PrevRange() // arbitrary, but somewhere inside the body means better diagnostics

			p.recover(end) // attempt to recover to the token after the end of this body
			break Token
		}
	}

	return &Body{
		Attributes: attrs,
		Blocks:     blocks,

		SrcRange: hcl.RangeBetween(startRange, endRange),
		EndRange: hcl.Range{
			Filename: endRange.Filename,
			Start:    endRange.End,
			End:      endRange.End,
		},
	}, diags
}

func (p *parser) ParseBodyItem() (Node, hcl.Diagnostics) {
	ident := p.Read()
	if ident.Type != TokenIdent {
		p.recoverAfterBodyItem()
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Attribute or block definition required",
				Detail:   "An attribute or block definition is required here.",
				Subject:  &ident.Range,
			},
		}
	}

	next := p.Peek()

	switch next.Type {
	case TokenEqual:
		return p.finishParsingBodyAttribute(ident)
	case TokenOQuote, TokenOBrace:
		return p.finishParsingBodyBlock(ident)
	default:
		p.recoverAfterBodyItem()
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Attribute or block definition required",
				Detail:   "An attribute or block definition is required here. To define an attribute, use the equals sign \"=\" to introduce the attribute value.",
				Subject:  &ident.Range,
			},
		}
	}

	return nil, nil
}

func (p *parser) finishParsingBodyAttribute(ident Token) (Node, hcl.Diagnostics) {
	eqTok := p.Read() // eat equals token
	if eqTok.Type != TokenEqual {
		// should never happen if caller behaves
		panic("finishParsingBodyAttribute called with next not equals")
	}

	var endRange hcl.Range

	expr, diags := p.ParseExpression()
	if p.recovery && diags.HasErrors() {
		// recovery within expressions tends to be tricky, so we've probably
		// landed somewhere weird. We'll try to reset to the start of a body
		// item so parsing can continue.
		endRange = p.PrevRange()
		p.recoverAfterBodyItem()
	} else {
		end := p.Peek()
		if end.Type != TokenNewline {
			if !p.recovery {
				if end.Type == TokenEOF {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing newline after attribute definition",
						Detail:   "A newline is required after an attribute definition at the end of a file.",
						Subject:  &end.Range,
						Context:  hcl.RangeBetween(ident.Range, end.Range).Ptr(),
					})
				} else {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing newline after attribute definition",
						Detail:   "An attribute definition must end with a newline.",
						Subject:  &end.Range,
						Context:  hcl.RangeBetween(ident.Range, end.Range).Ptr(),
					})
				}
			}
			endRange = p.PrevRange()
			p.recoverAfterBodyItem()
		} else {
			endRange = p.PrevRange()
			p.Read() // eat newline
		}
	}

	return &Attribute{
		Name: string(ident.Bytes),
		Expr: expr,

		SrcRange:    hcl.RangeBetween(ident.Range, endRange),
		NameRange:   ident.Range,
		EqualsRange: eqTok.Range,
	}, diags
}

func (p *parser) finishParsingBodyBlock(ident Token) (Node, hcl.Diagnostics) {
	var blockType = string(ident.Bytes)
	var diags hcl.Diagnostics
	var labels []string
	var labelRanges []hcl.Range

	var oBrace Token

Token:
	for {
		tok := p.Peek()

		switch tok.Type {

		case TokenOBrace:
			oBrace = p.Read()
			break Token

		case TokenOQuote:
			label, labelRange, labelDiags := p.parseQuotedStringLiteral()
			diags = append(diags, labelDiags...)
			labels = append(labels, label)
			labelRanges = append(labelRanges, labelRange)
			if labelDiags.HasErrors() {
				p.recoverAfterBodyItem()
				return &Block{
					Type:   blockType,
					Labels: labels,
					Body:   nil,

					TypeRange:       ident.Range,
					LabelRanges:     labelRanges,
					OpenBraceRange:  ident.Range, // placeholder
					CloseBraceRange: ident.Range, // placeholder
				}, diags
			}

		default:
			switch tok.Type {
			case TokenEqual:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid block definition",
					Detail:   "The equals sign \"=\" indicates an attribute definition, and must not be used when defining a block.",
					Subject:  &tok.Range,
					Context:  hcl.RangeBetween(ident.Range, tok.Range).Ptr(),
				})
			case TokenNewline:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid block definition",
					Detail:   "A block definition must have block content delimited by \"{\" and \"}\", starting on the same line as the block header.",
					Subject:  &tok.Range,
					Context:  hcl.RangeBetween(ident.Range, tok.Range).Ptr(),
				})
			default:
				if !p.recovery {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid block definition",
						Detail:   "Either a quoted string block label or an opening brace (\"{\") is expected here.",
						Subject:  &tok.Range,
						Context:  hcl.RangeBetween(ident.Range, tok.Range).Ptr(),
					})
				}
			}

			p.recoverAfterBodyItem()

			return &Block{
				Type:   blockType,
				Labels: labels,
				Body:   nil,

				TypeRange:       ident.Range,
				LabelRanges:     labelRanges,
				OpenBraceRange:  ident.Range, // placeholder
				CloseBraceRange: ident.Range, // placeholder
			}, diags
		}
	}

	// Once we fall out here, the peeker is pointed just after our opening
	// brace, so we can begin our nested body parsing.
	body, bodyDiags := p.ParseBody(TokenCBrace)
	diags = append(diags, bodyDiags...)
	cBraceRange := p.PrevRange()

	eol := p.Peek()
	if eol.Type == TokenNewline {
		p.Read() // eat newline
	} else {
		if !p.recovery {
			if eol.Type == TokenEOF {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing newline after block definition",
					Detail:   "A newline is required after a block definition at the end of a file.",
					Subject:  &eol.Range,
					Context:  hcl.RangeBetween(ident.Range, eol.Range).Ptr(),
				})
			} else {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing newline after block definition",
					Detail:   "A block definition must end with a newline.",
					Subject:  &eol.Range,
					Context:  hcl.RangeBetween(ident.Range, eol.Range).Ptr(),
				})
			}
		}
		p.recoverAfterBodyItem()
	}

	return &Block{
		Type:   blockType,
		Labels: labels,
		Body:   body,

		TypeRange:       ident.Range,
		LabelRanges:     labelRanges,
		OpenBraceRange:  oBrace.Range,
		CloseBraceRange: cBraceRange,
	}, diags
}

func (p *parser) ParseExpression() (Expression, hcl.Diagnostics) {
	return p.parseTernaryConditional()
}

func (p *parser) parseTernaryConditional() (Expression, hcl.Diagnostics) {
	// The ternary conditional operator (.. ? .. : ..) behaves somewhat
	// like a binary operator except that the "symbol" is itself
	// an expression enclosed in two punctuation characters.
	// The middle expression is parsed as if the ? and : symbols
	// were parentheses. The "rhs" (the "false expression") is then
	// treated right-associatively so it behaves similarly to the
	// middle in terms of precedence.

	startRange := p.NextRange()
	var condExpr, trueExpr, falseExpr Expression
	var diags hcl.Diagnostics

	condExpr, condDiags := p.parseBinaryOps(binaryOps)
	diags = append(diags, condDiags...)
	if p.recovery && condDiags.HasErrors() {
		return condExpr, diags
	}

	questionMark := p.Peek()
	if questionMark.Type != TokenQuestion {
		return condExpr, diags
	}

	p.Read() // eat question mark

	trueExpr, trueDiags := p.ParseExpression()
	diags = append(diags, trueDiags...)
	if p.recovery && trueDiags.HasErrors() {
		return condExpr, diags
	}

	colon := p.Peek()
	if colon.Type != TokenColon {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing false expression in conditional",
			Detail:   "The conditional operator (...?...:...) requires a false expression, delimited by a colon.",
			Subject:  &colon.Range,
			Context:  hcl.RangeBetween(startRange, colon.Range).Ptr(),
		})
		return condExpr, diags
	}

	p.Read() // eat colon

	falseExpr, falseDiags := p.ParseExpression()
	diags = append(diags, falseDiags...)
	if p.recovery && falseDiags.HasErrors() {
		return condExpr, diags
	}

	return &ConditionalExpr{
		Condition:   condExpr,
		TrueResult:  trueExpr,
		FalseResult: falseExpr,

		SrcRange: hcl.RangeBetween(startRange, falseExpr.Range()),
	}, diags
}

// parseBinaryOps calls itself recursively to work through all of the
// operator precedence groups, and then eventually calls parseExpressionTerm
// for each operand.
func (p *parser) parseBinaryOps(ops []map[TokenType]*Operation) (Expression, hcl.Diagnostics) {
	if len(ops) == 0 {
		// We've run out of operators, so now we'll just try to parse a term.
		return p.parseExpressionWithTraversals()
	}

	thisLevel := ops[0]
	remaining := ops[1:]

	var lhs, rhs Expression
	var operation *Operation
	var diags hcl.Diagnostics

	// Parse a term that might be the first operand of a binary
	// operation or it might just be a standalone term.
	// We won't know until we've parsed it and can look ahead
	// to see if there's an operator token for this level.
	lhs, lhsDiags := p.parseBinaryOps(remaining)
	diags = append(diags, lhsDiags...)
	if p.recovery && lhsDiags.HasErrors() {
		return lhs, diags
	}

	// We'll keep eating up operators until we run out, so that operators
	// with the same precedence will combine in a left-associative manner:
	// a+b+c => (a+b)+c, not a+(b+c)
	//
	// Should we later want to have right-associative operators, a way
	// to achieve that would be to call back up to ParseExpression here
	// instead of iteratively parsing only the remaining operators.
	for {
		next := p.Peek()
		var newOp *Operation
		var ok bool
		if newOp, ok = thisLevel[next.Type]; !ok {
			break
		}

		// Are we extending an expression started on the previous iteration?
		if operation != nil {
			lhs = &BinaryOpExpr{
				LHS: lhs,
				Op:  operation,
				RHS: rhs,

				SrcRange: hcl.RangeBetween(lhs.Range(), rhs.Range()),
			}
		}

		operation = newOp
		p.Read() // eat operator token
		var rhsDiags hcl.Diagnostics
		rhs, rhsDiags = p.parseBinaryOps(remaining)
		diags = append(diags, rhsDiags...)
		if p.recovery && rhsDiags.HasErrors() {
			return lhs, diags
		}
	}

	if operation == nil {
		return lhs, diags
	}

	return &BinaryOpExpr{
		LHS: lhs,
		Op:  operation,
		RHS: rhs,

		SrcRange: hcl.RangeBetween(lhs.Range(), rhs.Range()),
	}, diags
}

func (p *parser) parseExpressionWithTraversals() (Expression, hcl.Diagnostics) {
	term, diags := p.parseExpressionTerm()
	ret := term

Traversal:
	for {
		next := p.Peek()

		switch next.Type {
		case TokenDot:
			// Attribute access or splat
			dot := p.Read()
			attrTok := p.Peek()

			switch attrTok.Type {
			case TokenIdent:
				attrTok = p.Read() // eat token
				name := string(attrTok.Bytes)
				rng := hcl.RangeBetween(dot.Range, attrTok.Range)
				step := hcl.TraverseAttr{
					Name:     name,
					SrcRange: rng,
				}

				ret = makeRelativeTraversal(ret, step, rng)

			case TokenStar:
				// "Attribute-only" splat expression.
				// (This is a kinda weird construct inherited from HIL, which
				// behaves a bit like a [*] splat except that it is only able
				// to do attribute traversals into each of its elements,
				// whereas foo[*] can support _any_ traversal.
				marker := p.Read() // eat star
				trav := make(hcl.Traversal, 0, 1)
				var firstRange, lastRange hcl.Range
				firstRange = p.NextRange()
				for p.Peek().Type == TokenDot {
					dot := p.Read()

					if p.Peek().Type == TokenNumberLit {
						// Continuing the "weird stuff inherited from HIL"
						// theme, we also allow numbers as attribute names
						// inside splats and interpret them as indexing
						// into a list, for expressions like:
						// foo.bar.*.baz.0.foo
						numTok := p.Read()
						numVal, numDiags := p.numberLitValue(numTok)
						diags = append(diags, numDiags...)
						trav = append(trav, hcl.TraverseIndex{
							Key:      numVal,
							SrcRange: hcl.RangeBetween(dot.Range, numTok.Range),
						})
						lastRange = numTok.Range
						continue
					}

					if p.Peek().Type != TokenIdent {
						if !p.recovery {
							if p.Peek().Type == TokenStar {
								diags = append(diags, &hcl.Diagnostic{
									Severity: hcl.DiagError,
									Summary:  "Nested splat expression not allowed",
									Detail:   "A splat expression (*) cannot be used inside another attribute-only splat expression.",
									Subject:  p.Peek().Range.Ptr(),
								})
							} else {
								diags = append(diags, &hcl.Diagnostic{
									Severity: hcl.DiagError,
									Summary:  "Invalid attribute name",
									Detail:   "An attribute name is required after a dot.",
									Subject:  &attrTok.Range,
								})
							}
						}
						p.setRecovery()
						continue Traversal
					}

					attrTok := p.Read()
					trav = append(trav, hcl.TraverseAttr{
						Name:     string(attrTok.Bytes),
						SrcRange: hcl.RangeBetween(dot.Range, attrTok.Range),
					})
					lastRange = attrTok.Range
				}

				itemExpr := &AnonSymbolExpr{
					SrcRange: hcl.RangeBetween(dot.Range, marker.Range),
				}
				var travExpr Expression
				if len(trav) == 0 {
					travExpr = itemExpr
				} else {
					travExpr = &RelativeTraversalExpr{
						Source:    itemExpr,
						Traversal: trav,
						SrcRange:  hcl.RangeBetween(firstRange, lastRange),
					}
				}

				ret = &SplatExpr{
					Source: ret,
					Each:   travExpr,
					Item:   itemExpr,

					SrcRange:    hcl.RangeBetween(dot.Range, lastRange),
					MarkerRange: hcl.RangeBetween(dot.Range, marker.Range),
				}

			default:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid attribute name",
					Detail:   "An attribute name is required after a dot.",
					Subject:  &attrTok.Range,
				})
				// This leaves the peeker in a bad place, so following items
				// will probably be misparsed until we hit something that
				// allows us to re-sync.
				//
				// We will probably need to do something better here eventually
				// in order to support autocomplete triggered by typing a
				// period.
				p.setRecovery()
			}

		case TokenOBrack:
			// Indexing of a collection.
			// This may or may not be a hcl.Traverser, depending on whether
			// the key value is something constant.

			open := p.Read()
			// TODO: If we have a TokenStar inside our brackets, parse as
			// a Splat expression: foo[*].baz[0].
			var close Token
			p.PushIncludeNewlines(false) // arbitrary newlines allowed in brackets
			keyExpr, keyDiags := p.ParseExpression()
			diags = append(diags, keyDiags...)
			if p.recovery && keyDiags.HasErrors() {
				close = p.recover(TokenCBrack)
			} else {
				close = p.Read()
				if close.Type != TokenCBrack && !p.recovery {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing close bracket on index",
						Detail:   "The index operator must end with a closing bracket (\"]\").",
						Subject:  &close.Range,
					})
					close = p.recover(TokenCBrack)
				}
			}
			p.PushIncludeNewlines(true)

			if lit, isLit := keyExpr.(*LiteralValueExpr); isLit {
				litKey, _ := lit.Value(nil)
				rng := hcl.RangeBetween(open.Range, close.Range)
				step := &hcl.TraverseIndex{
					Key:      litKey,
					SrcRange: rng,
				}
				ret = makeRelativeTraversal(ret, step, rng)
			} else {
				rng := hcl.RangeBetween(open.Range, close.Range)
				ret = &IndexExpr{
					Collection: ret,
					Key:        keyExpr,

					SrcRange:  rng,
					OpenRange: open.Range,
				}
			}

		default:
			break Traversal
		}
	}

	return ret, diags
}

// makeRelativeTraversal takes an expression and a traverser and returns
// a traversal expression that combines the two. If the given expression
// is already a traversal, it is extended in place (mutating it) and
// returned. If it isn't, a new RelativeTraversalExpr is created and returned.
func makeRelativeTraversal(expr Expression, next hcl.Traverser, rng hcl.Range) Expression {
	switch texpr := expr.(type) {
	case *ScopeTraversalExpr:
		texpr.Traversal = append(texpr.Traversal, next)
		texpr.SrcRange = hcl.RangeBetween(texpr.SrcRange, rng)
		return texpr
	case *RelativeTraversalExpr:
		texpr.Traversal = append(texpr.Traversal, next)
		texpr.SrcRange = hcl.RangeBetween(texpr.SrcRange, rng)
		return texpr
	default:
		return &RelativeTraversalExpr{
			Source:    expr,
			Traversal: hcl.Traversal{next},
			SrcRange:  rng,
		}
	}
}

func (p *parser) parseExpressionTerm() (Expression, hcl.Diagnostics) {
	start := p.Peek()

	switch start.Type {
	case TokenOParen:
		p.Read() // eat open paren

		p.PushIncludeNewlines(false)

		expr, diags := p.ParseExpression()
		if diags.HasErrors() {
			// attempt to place the peeker after our closing paren
			// before we return, so that the next parser has some
			// chance of finding a valid expression.
			p.recover(TokenCParen)
			p.PopIncludeNewlines()
			return expr, diags
		}

		close := p.Peek()
		if close.Type != TokenCParen {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unbalanced parentheses",
				Detail:   "Expected a closing parenthesis to terminate the expression.",
				Subject:  &close.Range,
				Context:  hcl.RangeBetween(start.Range, close.Range).Ptr(),
			})
			p.setRecovery()
		}

		p.Read() // eat closing paren
		p.PopIncludeNewlines()

		return expr, diags

	case TokenNumberLit:
		tok := p.Read() // eat number token

		numVal, diags := p.numberLitValue(tok)
		return &LiteralValueExpr{
			Val:      numVal,
			SrcRange: tok.Range,
		}, diags

	case TokenIdent:
		tok := p.Read() // eat identifier token

		if p.Peek().Type == TokenOParen {
			return p.finishParsingFunctionCall(tok)
		}

		name := string(tok.Bytes)
		switch name {
		case "true":
			return &LiteralValueExpr{
				Val:      cty.True,
				SrcRange: tok.Range,
			}, nil
		case "false":
			return &LiteralValueExpr{
				Val:      cty.False,
				SrcRange: tok.Range,
			}, nil
		case "null":
			return &LiteralValueExpr{
				Val:      cty.NullVal(cty.DynamicPseudoType),
				SrcRange: tok.Range,
			}, nil
		default:
			return &ScopeTraversalExpr{
				Traversal: hcl.Traversal{
					hcl.TraverseRoot{
						Name:     name,
						SrcRange: tok.Range,
					},
				},
				SrcRange: tok.Range,
			}, nil
		}

	case TokenOQuote, TokenOHeredoc:
		open := p.Read() // eat opening marker
		closer := p.oppositeBracket(open.Type)
		exprs, passthru, _, diags := p.parseTemplateInner(closer)

		closeRange := p.PrevRange()

		if passthru {
			if len(exprs) != 1 {
				panic("passthru set with len(exprs) != 1")
			}
			return &TemplateWrapExpr{
				Wrapped:  exprs[0],
				SrcRange: hcl.RangeBetween(open.Range, closeRange),
			}, diags
		}

		return &TemplateExpr{
			Parts:    exprs,
			SrcRange: hcl.RangeBetween(open.Range, closeRange),
		}, diags

	case TokenMinus:
		tok := p.Read() // eat minus token

		// Important to use parseExpressionWithTraversals rather than parseExpression
		// here, otherwise we can capture a following binary expression into
		// our negation.
		// e.g. -46+5 should parse as (-46)+5, not -(46+5)
		operand, diags := p.parseExpressionWithTraversals()
		return &UnaryOpExpr{
			Op:  OpNegate,
			Val: operand,

			SrcRange:    hcl.RangeBetween(tok.Range, operand.Range()),
			SymbolRange: tok.Range,
		}, diags

	case TokenBang:
		tok := p.Read() // eat bang token

		// Important to use parseExpressionWithTraversals rather than parseExpression
		// here, otherwise we can capture a following binary expression into
		// our negation.
		operand, diags := p.parseExpressionWithTraversals()
		return &UnaryOpExpr{
			Op:  OpLogicalNot,
			Val: operand,

			SrcRange:    hcl.RangeBetween(tok.Range, operand.Range()),
			SymbolRange: tok.Range,
		}, diags

	case TokenOBrack:
		return p.parseTupleCons()

	case TokenOBrace:
		return p.parseObjectCons()

	default:
		var diags hcl.Diagnostics
		if !p.recovery {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid expression",
				Detail:   "Expected the start of an expression, but found an invalid expression token.",
				Subject:  &start.Range,
			})
		}
		p.setRecovery()

		// Return a placeholder so that the AST is still structurally sound
		// even in the presence of parse errors.
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: start.Range,
		}, diags
	}
}

func (p *parser) numberLitValue(tok Token) (cty.Value, hcl.Diagnostics) {
	// We'll lean on the cty converter to do the conversion, to ensure that
	// the behavior is the same as what would happen if converting a
	// non-literal string to a number.
	numStrVal := cty.StringVal(string(tok.Bytes))
	numVal, err := convert.Convert(numStrVal, cty.Number)
	if err != nil {
		ret := cty.UnknownVal(cty.Number)
		return ret, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid number literal",
				// FIXME: not a very good error message, but convert only
				// gives us "a number is required", so not much help either.
				Detail:  "Failed to recognize the value of this number literal.",
				Subject: &tok.Range,
			},
		}
	}
	return numVal, nil
}

// finishParsingFunctionCall parses a function call assuming that the function
// name was already read, and so the peeker should be pointing at the opening
// parenthesis after the name.
func (p *parser) finishParsingFunctionCall(name Token) (Expression, hcl.Diagnostics) {
	openTok := p.Read()
	if openTok.Type != TokenOParen {
		// should never happen if callers behave
		panic("finishParsingFunctionCall called with non-parenthesis as next token")
	}

	var args []Expression
	var diags hcl.Diagnostics
	var expandFinal bool
	var closeTok Token

	// Arbitrary newlines are allowed inside the function call parentheses.
	p.PushIncludeNewlines(false)

Token:
	for {
		tok := p.Peek()

		if tok.Type == TokenCParen {
			closeTok = p.Read() // eat closing paren
			break Token
		}

		arg, argDiags := p.ParseExpression()
		args = append(args, arg)
		diags = append(diags, argDiags...)
		if p.recovery && argDiags.HasErrors() {
			// if there was a parse error in the argument then we've
			// probably been left in a weird place in the token stream,
			// so we'll bail out with a partial argument list.
			p.recover(TokenCParen)
			break Token
		}

		sep := p.Read()
		if sep.Type == TokenCParen {
			closeTok = sep
			break Token
		}

		if sep.Type == TokenEllipsis {
			expandFinal = true

			if p.Peek().Type != TokenCParen {
				if !p.recovery {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing closing parenthesis",
						Detail:   "An expanded function argument (with ...) must be immediately followed by closing parentheses.",
						Subject:  &sep.Range,
						Context:  hcl.RangeBetween(name.Range, sep.Range).Ptr(),
					})
				}
				closeTok = p.recover(TokenCParen)
			} else {
				closeTok = p.Read() // eat closing paren
			}
			break Token
		}

		if sep.Type != TokenComma {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing argument separator",
				Detail:   "A comma is required to separate each function argument from the next.",
				Subject:  &sep.Range,
				Context:  hcl.RangeBetween(name.Range, sep.Range).Ptr(),
			})
			closeTok = p.recover(TokenCParen)
			break Token
		}

		if p.Peek().Type == TokenCParen {
			// A trailing comma after the last argument gets us in here.
			closeTok = p.Read() // eat closing paren
			break Token
		}

	}

	p.PopIncludeNewlines()

	return &FunctionCallExpr{
		Name: string(name.Bytes),
		Args: args,

		ExpandFinal: expandFinal,

		NameRange:       name.Range,
		OpenParenRange:  openTok.Range,
		CloseParenRange: closeTok.Range,
	}, diags
}

func (p *parser) parseTupleCons() (Expression, hcl.Diagnostics) {
	open := p.Read()
	if open.Type != TokenOBrack {
		// Should never happen if callers are behaving
		panic("parseTupleCons called without peeker pointing to open bracket")
	}

	p.PushIncludeNewlines(false)
	defer p.PopIncludeNewlines()

	if forKeyword.TokenMatches(p.Peek()) {
		return p.finishParsingForExpr(open)
	}

	var close Token

	var diags hcl.Diagnostics
	var exprs []Expression

	for {
		next := p.Peek()
		if next.Type == TokenCBrack {
			close = p.Read() // eat closer
			break
		}

		expr, exprDiags := p.ParseExpression()
		exprs = append(exprs, expr)
		diags = append(diags, exprDiags...)

		if p.recovery && exprDiags.HasErrors() {
			// If expression parsing failed then we are probably in a strange
			// place in the token stream, so we'll bail out and try to reset
			// to after our closing bracket to allow parsing to continue.
			close = p.recover(TokenCBrack)
			break
		}

		next = p.Peek()
		if next.Type == TokenCBrack {
			close = p.Read() // eat closer
			break
		}

		if next.Type != TokenComma {
			if !p.recovery {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing item separator",
					Detail:   "Expected a comma to mark the beginning of the next item.",
					Subject:  &next.Range,
					Context:  hcl.RangeBetween(open.Range, next.Range).Ptr(),
				})
			}
			close = p.recover(TokenCBrack)
			break
		}

		p.Read() // eat comma

	}

	return &TupleConsExpr{
		Exprs: exprs,

		SrcRange:  hcl.RangeBetween(open.Range, close.Range),
		OpenRange: open.Range,
	}, diags
}

func (p *parser) parseObjectCons() (Expression, hcl.Diagnostics) {
	open := p.Read()
	if open.Type != TokenOBrace {
		// Should never happen if callers are behaving
		panic("parseObjectCons called without peeker pointing to open brace")
	}

	p.PushIncludeNewlines(true)
	defer p.PopIncludeNewlines()

	if forKeyword.TokenMatches(p.Peek()) {
		return p.finishParsingForExpr(open)
	}

	var close Token

	var diags hcl.Diagnostics
	var items []ObjectConsItem

	for {
		next := p.Peek()
		if next.Type == TokenNewline {
			p.Read() // eat newline
			continue
		}

		if next.Type == TokenCBrace {
			close = p.Read() // eat closer
			break
		}

		// As a special case, we allow the key to be a literal identifier.
		// This means that a variable reference or function call can't appear
		// directly as key expression, and must instead be wrapped in some
		// disambiguation punctuation, like (var.a) = "b" or "${var.a}" = "b".
		var key Expression
		var keyDiags hcl.Diagnostics
		if p.Peek().Type == TokenIdent {
			nameTok := p.Read()
			key = &LiteralValueExpr{
				Val: cty.StringVal(string(nameTok.Bytes)),

				SrcRange: nameTok.Range,
			}
		} else {
			key, keyDiags = p.ParseExpression()
		}

		diags = append(diags, keyDiags...)

		if p.recovery && keyDiags.HasErrors() {
			// If expression parsing failed then we are probably in a strange
			// place in the token stream, so we'll bail out and try to reset
			// to after our closing brace to allow parsing to continue.
			close = p.recover(TokenCBrace)
			break
		}

		next = p.Peek()
		if next.Type != TokenEqual && next.Type != TokenColon {
			if !p.recovery {
				if next.Type == TokenNewline || next.Type == TokenComma {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing item value",
						Detail:   "Expected an item value, introduced by an equals sign (\"=\").",
						Subject:  &next.Range,
						Context:  hcl.RangeBetween(open.Range, next.Range).Ptr(),
					})
				} else {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Missing key/value separator",
						Detail:   "Expected an equals sign (\"=\") to mark the beginning of the item value.",
						Subject:  &next.Range,
						Context:  hcl.RangeBetween(open.Range, next.Range).Ptr(),
					})
				}
			}
			close = p.recover(TokenCBrace)
			break
		}

		p.Read() // eat equals sign or colon

		value, valueDiags := p.ParseExpression()
		diags = append(diags, valueDiags...)

		if p.recovery && valueDiags.HasErrors() {
			// If expression parsing failed then we are probably in a strange
			// place in the token stream, so we'll bail out and try to reset
			// to after our closing brace to allow parsing to continue.
			close = p.recover(TokenCBrace)
			break
		}

		items = append(items, ObjectConsItem{
			KeyExpr:   key,
			ValueExpr: value,
		})

		next = p.Peek()
		if next.Type == TokenCBrace {
			close = p.Read() // eat closer
			break
		}

		if next.Type != TokenComma && next.Type != TokenNewline {
			if !p.recovery {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing item separator",
					Detail:   "Expected a newline or comma to mark the beginning of the next item.",
					Subject:  &next.Range,
					Context:  hcl.RangeBetween(open.Range, next.Range).Ptr(),
				})
			}
			close = p.recover(TokenCBrace)
			break
		}

		p.Read() // eat comma or newline

	}

	return &ObjectConsExpr{
		Items: items,

		SrcRange:  hcl.RangeBetween(open.Range, close.Range),
		OpenRange: open.Range,
	}, diags
}

func (p *parser) finishParsingForExpr(open Token) (Expression, hcl.Diagnostics) {
	introducer := p.Read()
	if !forKeyword.TokenMatches(introducer) {
		// Should never happen if callers are behaving
		panic("finishParsingForExpr called without peeker pointing to 'for' identifier")
	}

	var makeObj bool
	var closeType TokenType
	switch open.Type {
	case TokenOBrace:
		makeObj = true
		closeType = TokenCBrace
	case TokenOBrack:
		makeObj = false // making a tuple
		closeType = TokenCBrack
	default:
		// Should never happen if callers are behaving
		panic("finishParsingForExpr called with invalid open token")
	}

	var diags hcl.Diagnostics
	var keyName, valName string

	if p.Peek().Type != TokenIdent {
		if !p.recovery {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "For expression requires variable name after 'for'.",
				Subject:  p.Peek().Range.Ptr(),
				Context:  hcl.RangeBetween(open.Range, p.Peek().Range).Ptr(),
			})
		}
		close := p.recover(closeType)
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: hcl.RangeBetween(open.Range, close.Range),
		}, diags
	}

	valName = string(p.Read().Bytes)

	if p.Peek().Type == TokenComma {
		// What we just read was actually the key, then.
		keyName = valName
		p.Read() // eat comma

		if p.Peek().Type != TokenIdent {
			if !p.recovery {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid 'for' expression",
					Detail:   "For expression requires value variable name after comma.",
					Subject:  p.Peek().Range.Ptr(),
					Context:  hcl.RangeBetween(open.Range, p.Peek().Range).Ptr(),
				})
			}
			close := p.recover(closeType)
			return &LiteralValueExpr{
				Val:      cty.DynamicVal,
				SrcRange: hcl.RangeBetween(open.Range, close.Range),
			}, diags
		}

		valName = string(p.Read().Bytes)
	}

	if !inKeyword.TokenMatches(p.Peek()) {
		if !p.recovery {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "For expression requires 'in' keyword after names.",
				Subject:  p.Peek().Range.Ptr(),
				Context:  hcl.RangeBetween(open.Range, p.Peek().Range).Ptr(),
			})
		}
		close := p.recover(closeType)
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: hcl.RangeBetween(open.Range, close.Range),
		}, diags
	}
	p.Read() // eat 'in' keyword

	collExpr, collDiags := p.ParseExpression()
	diags = append(diags, collDiags...)
	if p.recovery && collDiags.HasErrors() {
		close := p.recover(closeType)
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: hcl.RangeBetween(open.Range, close.Range),
		}, diags
	}

	if p.Peek().Type != TokenColon {
		if !p.recovery {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "For expression requires colon after collection expression.",
				Subject:  p.Peek().Range.Ptr(),
				Context:  hcl.RangeBetween(open.Range, p.Peek().Range).Ptr(),
			})
		}
		close := p.recover(closeType)
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: hcl.RangeBetween(open.Range, close.Range),
		}, diags
	}
	p.Read() // eat colon

	var keyExpr, valExpr Expression
	var keyDiags, valDiags hcl.Diagnostics
	valExpr, valDiags = p.ParseExpression()
	if p.Peek().Type == TokenFatArrow {
		// What we just parsed was actually keyExpr
		p.Read() // eat the fat arrow
		keyExpr, keyDiags = valExpr, valDiags

		valExpr, valDiags = p.ParseExpression()
	}
	diags = append(diags, keyDiags...)
	diags = append(diags, valDiags...)
	if p.recovery && (keyDiags.HasErrors() || valDiags.HasErrors()) {
		close := p.recover(closeType)
		return &LiteralValueExpr{
			Val:      cty.DynamicVal,
			SrcRange: hcl.RangeBetween(open.Range, close.Range),
		}, diags
	}

	group := false
	var ellipsis Token
	if p.Peek().Type == TokenEllipsis {
		ellipsis = p.Read()
		group = true
	}

	var condExpr Expression
	var condDiags hcl.Diagnostics
	if ifKeyword.TokenMatches(p.Peek()) {
		p.Read() // eat "if"
		condExpr, condDiags = p.ParseExpression()
		diags = append(diags, condDiags...)
		if p.recovery && condDiags.HasErrors() {
			close := p.recover(p.oppositeBracket(open.Type))
			return &LiteralValueExpr{
				Val:      cty.DynamicVal,
				SrcRange: hcl.RangeBetween(open.Range, close.Range),
			}, diags
		}
	}

	var close Token
	if p.Peek().Type == closeType {
		close = p.Read()
	} else {
		if !p.recovery {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "Extra characters after the end of the 'for' expression.",
				Subject:  p.Peek().Range.Ptr(),
				Context:  hcl.RangeBetween(open.Range, p.Peek().Range).Ptr(),
			})
		}
		close = p.recover(closeType)
	}

	if !makeObj {
		if keyExpr != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "Key expression is not valid when building a tuple.",
				Subject:  keyExpr.Range().Ptr(),
				Context:  hcl.RangeBetween(open.Range, close.Range).Ptr(),
			})
		}

		if group {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "Grouping ellipsis (...) cannot be used when building a tuple.",
				Subject:  &ellipsis.Range,
				Context:  hcl.RangeBetween(open.Range, close.Range).Ptr(),
			})
		}
	} else {
		if keyExpr == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid 'for' expression",
				Detail:   "Key expression is required when building an object.",
				Subject:  valExpr.Range().Ptr(),
				Context:  hcl.RangeBetween(open.Range, close.Range).Ptr(),
			})
		}
	}

	return &ForExpr{
		KeyVar:   keyName,
		ValVar:   valName,
		CollExpr: collExpr,
		KeyExpr:  keyExpr,
		ValExpr:  valExpr,
		CondExpr: condExpr,
		Group:    group,

		SrcRange:   hcl.RangeBetween(open.Range, close.Range),
		OpenRange:  open.Range,
		CloseRange: close.Range,
	}, diags
}

// parseQuotedStringLiteral is a helper for parsing quoted strings that
// aren't allowed to contain any interpolations, such as block labels.
func (p *parser) parseQuotedStringLiteral() (string, hcl.Range, hcl.Diagnostics) {
	oQuote := p.Read()
	if oQuote.Type != TokenOQuote {
		return "", oQuote.Range, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid string literal",
				Detail:   "A quoted string is required here.",
				Subject:  &oQuote.Range,
			},
		}
	}

	var diags hcl.Diagnostics
	ret := &bytes.Buffer{}
	var cQuote Token

Token:
	for {
		tok := p.Read()
		switch tok.Type {

		case TokenCQuote:
			cQuote = tok
			break Token

		case TokenQuotedLit:
			s, sDiags := p.decodeStringLit(tok)
			diags = append(diags, sDiags...)
			ret.WriteString(s)

		case TokenTemplateControl, TokenTemplateInterp:
			which := "$"
			if tok.Type == TokenTemplateControl {
				which = "!"
			}

			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid string literal",
				Detail: fmt.Sprintf(
					"Template sequences are not allowed in this string. To include a literal %q, double it (as \"%s%s\") to escape it.",
					which, which, which,
				),
				Subject: &tok.Range,
				Context: hcl.RangeBetween(oQuote.Range, tok.Range).Ptr(),
			})
			p.recover(TokenTemplateSeqEnd)

		case TokenEOF:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unterminated string literal",
				Detail:   "Unable to find the closing quote mark before the end of the file.",
				Subject:  &tok.Range,
				Context:  hcl.RangeBetween(oQuote.Range, tok.Range).Ptr(),
			})
			break Token

		default:
			// Should never happen, as long as the scanner is behaving itself
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid string literal",
				Detail:   "This item is not valid in a string literal.",
				Subject:  &tok.Range,
				Context:  hcl.RangeBetween(oQuote.Range, tok.Range).Ptr(),
			})
			p.recover(TokenOQuote)
			break Token

		}

	}

	return ret.String(), hcl.RangeBetween(oQuote.Range, cQuote.Range), diags
}

// decodeStringLit processes the given token, which must be either a
// TokenQuotedLit or a TokenStringLit, returning the string resulting from
// resolving any escape sequences.
//
// If any error diagnostics are returned, the returned string may be incomplete
// or otherwise invalid.
func (p *parser) decodeStringLit(tok Token) (string, hcl.Diagnostics) {
	var quoted bool
	switch tok.Type {
	case TokenQuotedLit:
		quoted = true
	case TokenStringLit:
		quoted = false
	default:
		panic("decodeQuotedLit can only be used with TokenStringLit and TokenQuotedLit tokens")
	}
	var diags hcl.Diagnostics

	ret := make([]byte, 0, len(tok.Bytes))
	var esc []byte

	sc := bufio.NewScanner(bytes.NewReader(tok.Bytes))
	sc.Split(textseg.ScanGraphemeClusters)

	pos := tok.Range.Start
	newPos := pos
Character:
	for sc.Scan() {
		pos = newPos
		ch := sc.Bytes()

		// Adjust position based on our new character.
		// \r\n is considered to be a single character in text segmentation,
		if (len(ch) == 1 && ch[0] == '\n') || (len(ch) == 2 && ch[1] == '\n') {
			newPos.Line++
			newPos.Column = 0
		} else {
			newPos.Column++
		}
		newPos.Byte += len(ch)

		if len(esc) > 0 {
			switch esc[0] {
			case '\\':
				if len(ch) == 1 {
					switch ch[0] {

					// TODO: numeric character escapes with \uXXXX

					case 'n':
						ret = append(ret, '\n')
						esc = esc[:0]
						continue Character
					case 'r':
						ret = append(ret, '\r')
						esc = esc[:0]
						continue Character
					case 't':
						ret = append(ret, '\t')
						esc = esc[:0]
						continue Character
					case '"':
						ret = append(ret, '"')
						esc = esc[:0]
						continue Character
					case '\\':
						ret = append(ret, '\\')
						esc = esc[:0]
						continue Character
					}
				}

				var detail string
				switch {
				case len(ch) == 1 && (ch[0] == '$' || ch[0] == '!'):
					detail = fmt.Sprintf(
						"The characters \"\\%s\" do not form a recognized escape sequence. To escape a \"%s{\" template sequence, use \"%s%s{\".",
						ch, ch, ch, ch,
					)
				default:
					detail = fmt.Sprintf("The characters \"\\%s\" do not form a recognized escape sequence.", ch)
				}

				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid escape sequence",
					Detail:   detail,
					Subject: &hcl.Range{
						Filename: tok.Range.Filename,
						Start: hcl.Pos{
							Line:   pos.Line,
							Column: pos.Column - 1, // safe because we know the previous character must be a backslash
							Byte:   pos.Byte - 1,
						},
						End: hcl.Pos{
							Line:   pos.Line,
							Column: pos.Column + 1, // safe because we know the previous character must be a backslash
							Byte:   pos.Byte + len(ch),
						},
					},
				})
				ret = append(ret, ch...)
				esc = esc[:0]
				continue Character

			case '$', '!':
				switch len(esc) {
				case 1:
					if len(ch) == 1 && ch[0] == esc[0] {
						esc = append(esc, ch[0])
						continue Character
					}

					// Any other character means this wasn't an escape sequence
					// after all.
					ret = append(ret, esc...)
					ret = append(ret, ch...)
					esc = esc[:0]
				case 2:
					if len(ch) == 1 && ch[0] == '{' {
						// successful escape sequence
						ret = append(ret, esc[0])
					} else {
						// not an escape sequence, so just output literal
						ret = append(ret, esc...)
					}
					ret = append(ret, ch...)
					esc = esc[:0]
				default:
					// should never happen
					panic("have invalid escape sequence >2 characters")
				}

			}
		} else {
			if len(ch) == 1 {
				switch ch[0] {
				case '\\':
					if quoted { // ignore backslashes in unquoted mode
						esc = append(esc, '\\')
						continue Character
					}
				case '$':
					esc = append(esc, '$')
					continue Character
				case '!':
					esc = append(esc, '!')
					continue Character
				}
			}
			ret = append(ret, ch...)
		}
	}

	return string(ret), diags
}

// setRecovery turns on recovery mode without actually doing any recovery.
// This can be used when a parser knowingly leaves the peeker in a useless
// place and wants to suppress errors that might result from that decision.
func (p *parser) setRecovery() {
	p.recovery = true
}

// recover seeks forward in the token stream until it finds TokenType "end",
// then returns with the peeker pointed at the following token.
//
// If the given token type is a bracketer, this function will additionally
// count nested instances of the brackets to try to leave the peeker at
// the end of the _current_ instance of that bracketer, skipping over any
// nested instances. This is a best-effort operation and may have
// unpredictable results on input with bad bracketer nesting.
func (p *parser) recover(end TokenType) Token {
	start := p.oppositeBracket(end)
	p.recovery = true

	nest := 0
	for {
		tok := p.Read()
		ty := tok.Type
		if end == TokenTemplateSeqEnd && ty == TokenTemplateControl {
			// normalize so that our matching behavior can work, since
			// TokenTemplateControl/TokenTemplateInterp are asymmetrical
			// with TokenTemplateSeqEnd and thus we need to count both
			// openers if that's the closer we're looking for.
			ty = TokenTemplateInterp
		}

		switch ty {
		case start:
			nest++
		case end:
			if nest < 1 {
				return tok
			}

			nest--
		case TokenEOF:
			return tok
		}
	}
}

// recoverOver seeks forward in the token stream until it finds a block
// starting with TokenType "start", then finds the corresponding end token,
// leaving the peeker pointed at the token after that end token.
//
// The given token type _must_ be a bracketer. For example, if the given
// start token is TokenOBrace then the parser will be left at the _end_ of
// the next brace-delimited block encountered, or at EOF if no such block
// is found or it is unclosed.
func (p *parser) recoverOver(start TokenType) {
	end := p.oppositeBracket(start)

	// find the opening bracket first
Token:
	for {
		tok := p.Read()
		switch tok.Type {
		case start, TokenEOF:
			break Token
		}
	}

	// Now use our existing recover function to locate the _end_ of the
	// container we've found.
	p.recover(end)
}

func (p *parser) recoverAfterBodyItem() {
	p.recovery = true
	var open []TokenType

Token:
	for {
		tok := p.Read()

		switch tok.Type {

		case TokenNewline:
			if len(open) == 0 {
				break Token
			}

		case TokenEOF:
			break Token

		case TokenOBrace, TokenOBrack, TokenOParen, TokenOQuote, TokenOHeredoc, TokenTemplateInterp, TokenTemplateControl:
			open = append(open, tok.Type)

		case TokenCBrace, TokenCBrack, TokenCParen, TokenCQuote, TokenCHeredoc:
			opener := p.oppositeBracket(tok.Type)
			for len(open) > 0 && open[len(open)-1] != opener {
				open = open[:len(open)-1]
			}
			if len(open) > 0 {
				open = open[:len(open)-1]
			}

		case TokenTemplateSeqEnd:
			for len(open) > 0 && open[len(open)-1] != TokenTemplateInterp && open[len(open)-1] != TokenTemplateControl {
				open = open[:len(open)-1]
			}
			if len(open) > 0 {
				open = open[:len(open)-1]
			}

		}
	}
}

// oppositeBracket finds the bracket that opposes the given bracketer, or
// NilToken if the given token isn't a bracketer.
//
// "Bracketer", for the sake of this function, is one end of a matching
// open/close set of tokens that establish a bracketing context.
func (p *parser) oppositeBracket(ty TokenType) TokenType {
	switch ty {

	case TokenOBrace:
		return TokenCBrace
	case TokenOBrack:
		return TokenCBrack
	case TokenOParen:
		return TokenCParen
	case TokenOQuote:
		return TokenCQuote
	case TokenOHeredoc:
		return TokenCHeredoc

	case TokenCBrace:
		return TokenOBrace
	case TokenCBrack:
		return TokenOBrack
	case TokenCParen:
		return TokenOParen
	case TokenCQuote:
		return TokenOQuote
	case TokenCHeredoc:
		return TokenOHeredoc

	case TokenTemplateControl:
		return TokenTemplateSeqEnd
	case TokenTemplateInterp:
		return TokenTemplateSeqEnd
	case TokenTemplateSeqEnd:
		// This is ambigous, but we return Interp here because that's
		// what's assumed by the "recover" method.
		return TokenTemplateInterp

	default:
		return TokenNil
	}
}

func errPlaceholderExpr(rng hcl.Range) Expression {
	return &LiteralValueExpr{
		Val:      cty.DynamicVal,
		SrcRange: rng,
	}
}
