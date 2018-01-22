package zclsyntax

import (
    "bytes"

    "github.com/zclconf/go-zcl/zcl"
)

// This file is generated from scan_tokens.rl. DO NOT EDIT.
%%{
  # (except you are actually in scan_tokens.rl here, so edit away!)

  machine zcltok;
  write data;
}%%

func scanTokens(data []byte, filename string, start hcl.Pos, mode scanMode) []Token {
    f := &tokenAccum{
        Filename: filename,
        Bytes:    data,
        Pos:      start,
    }

    %%{
        include UnicodeDerived "unicode_derived.rl";

        UTF8Cont = 0x80 .. 0xBF;
        AnyUTF8 = (
            0x00..0x7F |
            0xC0..0xDF . UTF8Cont |
            0xE0..0xEF . UTF8Cont . UTF8Cont |
            0xF0..0xF7 . UTF8Cont . UTF8Cont . UTF8Cont
        );
        BrokenUTF8 = any - AnyUTF8;

        NumberLitContinue = (digit|'.'|('e'|'E') ('+'|'-')? digit);
        NumberLit = digit ("" | (NumberLitContinue - '.') | (NumberLitContinue* (NumberLitContinue - '.')));
        Ident = ID_Start (ID_Continue | '-')*;

        # Symbols that just represent themselves are handled as a single rule.
        SelfToken = "[" | "]" | "(" | ")" | "." | "," | "*" | "/" | "+" | "-" | "=" | "<" | ">" | "!" | "?" | ":" | "\n" | "&" | "|" | "~" | "^" | ";" | "`";

        EqualOp = "==";
        NotEqual = "!=";
        GreaterThanEqual = ">=";
        LessThanEqual = "<=";
        LogicalAnd = "&&";
        LogicalOr = "||";

        Ellipsis = "...";
        FatArrow = "=>";

        Newline = '\r' ? '\n';
        EndOfLine = Newline;

        BeginStringTmpl = '"';
        BeginHeredocTmpl = '<<' ('-')? Ident Newline;

        Comment = (
            ("#" (any - EndOfLine)* EndOfLine) |
            ("//" (any - EndOfLine)* EndOfLine) |
            ("/*" any* "*/")
        );

        # Tabs are not valid, but we accept them in the scanner and mark them
        # as tokens so that we can produce diagnostics advising the user to
        # use spaces instead.
        Tabs = 0x09+;

        # Note: zclwrite assumes that only ASCII spaces appear between tokens,
        # and uses this assumption to recreate the spaces between tokens by
        # looking at byte offset differences.
        Spaces = ' '+;

        action beginStringTemplate {
            token(TokenOQuote);
            fcall stringTemplate;
        }

        action endStringTemplate {
            token(TokenCQuote);
            fret;
        }

        action beginHeredocTemplate {
            token(TokenOHeredoc);
            // the token is currently the whole heredoc introducer, like
            // <<EOT or <<-EOT, followed by a newline. We want to extract
            // just the "EOT" portion that we'll use as the closing marker.

            marker := data[ts+2:te-1]
            if marker[0] == '-' {
                marker = marker[1:]
            }
            if marker[len(marker)-1] == '\r' {
                marker = marker[:len(marker)-1]
            }

            heredocs = append(heredocs, heredocInProgress{
                Marker:      marker,
                StartOfLine: true,
            })

            fcall heredocTemplate;
        }

        action heredocLiteralEOL {
            // This action is called specificially when a heredoc literal
            // ends with a newline character.

            // This might actually be our end marker.
            topdoc := &heredocs[len(heredocs)-1]
            if topdoc.StartOfLine {
                maybeMarker := bytes.TrimSpace(data[ts:te])
                if bytes.Equal(maybeMarker, topdoc.Marker) {
                    token(TokenCHeredoc);
                    heredocs = heredocs[:len(heredocs)-1]
                    fret;
                }
            }

            topdoc.StartOfLine = true;
            token(TokenStringLit);
        }

        action heredocLiteralMidline {
            // This action is called when a heredoc literal _doesn't_ end
            // with a newline character, e.g. because we're about to enter
            // an interpolation sequence.
            heredocs[len(heredocs)-1].StartOfLine = false;
            token(TokenStringLit);
        }

        action bareTemplateLiteral {
            token(TokenStringLit);
        }

        action beginTemplateInterp {
            token(TokenTemplateInterp);
            braces++;
            retBraces = append(retBraces, braces);
            if len(heredocs) > 0 {
                heredocs[len(heredocs)-1].StartOfLine = false;
            }
            fcall main;
        }

        action beginTemplateControl {
            token(TokenTemplateControl);
            braces++;
            retBraces = append(retBraces, braces);
            if len(heredocs) > 0 {
                heredocs[len(heredocs)-1].StartOfLine = false;
            }
            fcall main;
        }

        action openBrace {
            token(TokenOBrace);
            braces++;
        }

        action closeBrace {
            if len(retBraces) > 0 && retBraces[len(retBraces)-1] == braces {
                token(TokenTemplateSeqEnd);
                braces--;
                retBraces = retBraces[0:len(retBraces)-1]
                fret;
            } else {
                token(TokenCBrace);
                braces--;
            }
        }

        action closeTemplateSeqEatWhitespace {
            // Only consume from the retBraces stack and return if we are at
            // a suitable brace nesting level, otherwise things will get
            // confused. (Not entering this branch indicates a syntax error,
            // which we will catch in the parser.)
            if len(retBraces) > 0 && retBraces[len(retBraces)-1] == braces {
                token(TokenTemplateSeqEnd);
                braces--;
                retBraces = retBraces[0:len(retBraces)-1]
                fret;
            } else {
                // We intentionally generate a TokenTemplateSeqEnd here,
                // even though the user apparently wanted a brace, because
                // we want to allow the parser to catch the incorrect use
                // of a ~} to balance a generic opening brace, rather than
                // a template sequence.
                token(TokenTemplateSeqEnd);
                braces--;
            }
        }

        TemplateInterp = "${" ("~")?;
        TemplateControl = "%{" ("~")?;
        EndStringTmpl = '"';
        StringLiteralChars = (AnyUTF8 - ("\r"|"\n"));
        TemplateStringLiteral = (
            ('$' ^'{') |
            ('%' ^'{') |
            ('\\' StringLiteralChars) |
            (StringLiteralChars - ("$" | '%' | '"'))
        )+;
        HeredocStringLiteral = (
            ('$' ^'{') |
            ('%' ^'{') |
            (StringLiteralChars - ("$" | '%'))
        )*;
        BareStringLiteral = (
            ('$' ^'{') |
            ('%' ^'{') |
            (StringLiteralChars - ("$" | '%'))
        )* Newline?;

        stringTemplate := |*
            TemplateInterp        => beginTemplateInterp;
            TemplateControl       => beginTemplateControl;
            EndStringTmpl         => endStringTemplate;
            TemplateStringLiteral => { token(TokenQuotedLit); };
            AnyUTF8               => { token(TokenInvalid); };
            BrokenUTF8            => { token(TokenBadUTF8); };
        *|;

        heredocTemplate := |*
            TemplateInterp        => beginTemplateInterp;
            TemplateControl       => beginTemplateControl;
            HeredocStringLiteral EndOfLine => heredocLiteralEOL;
            HeredocStringLiteral  => heredocLiteralMidline;
            BrokenUTF8            => { token(TokenBadUTF8); };
        *|;

        bareTemplate := |*
            TemplateInterp        => beginTemplateInterp;
            TemplateControl       => beginTemplateControl;
            BareStringLiteral     => bareTemplateLiteral;
            BrokenUTF8            => { token(TokenBadUTF8); };
        *|;

        main := |*
            Spaces           => {};
            NumberLit        => { token(TokenNumberLit) };
            Ident            => { token(TokenIdent) };

            Comment          => { token(TokenComment) };
            Newline          => { token(TokenNewline) };

            EqualOp          => { token(TokenEqualOp); };
            NotEqual         => { token(TokenNotEqual); };
            GreaterThanEqual => { token(TokenGreaterThanEq); };
            LessThanEqual    => { token(TokenLessThanEq); };
            LogicalAnd       => { token(TokenAnd); };
            LogicalOr        => { token(TokenOr); };
            Ellipsis         => { token(TokenEllipsis); };
            FatArrow         => { token(TokenFatArrow); };
            SelfToken        => { selfToken() };

            "{"              => openBrace;
            "}"              => closeBrace;

            "~}"             => closeTemplateSeqEatWhitespace;

            BeginStringTmpl  => beginStringTemplate;
            BeginHeredocTmpl => beginHeredocTemplate;

            Tabs             => { token(TokenTabs) };
            BrokenUTF8       => { token(TokenBadUTF8) };
            AnyUTF8          => { token(TokenInvalid) };
        *|;

    }%%

    // Ragel state
	p := 0  // "Pointer" into data
	pe := len(data) // End-of-data "pointer"
    ts := 0
    te := 0
    act := 0
    eof := pe
    var stack []int
    var top int

    var cs int // current state
    switch mode {
    case scanNormal:
        cs = zcltok_en_main
    case scanTemplate:
        cs = zcltok_en_bareTemplate
    default:
        panic("invalid scanMode")
    }

    braces := 0
    var retBraces []int // stack of brace levels that cause us to use fret
    var heredocs []heredocInProgress // stack of heredocs we're currently processing

    %%{
        prepush {
            stack = append(stack, 0);
        }
        postpop {
            stack = stack[:len(stack)-1];
        }
    }%%

    // Make Go compiler happy
    _ = ts
    _ = te
    _ = act
    _ = eof

    token := func (ty TokenType) {
        f.emitToken(ty, ts, te)
    }
    selfToken := func () {
        b := data[ts:te]
        if len(b) != 1 {
            // should never happen
            panic("selfToken only works for single-character tokens")
        }
        f.emitToken(TokenType(b[0]), ts, te)
    }

    %%{
        write init nocs;
        write exec;
    }%%

    // If we fall out here without being in a final state then we've
    // encountered something that the scanner can't match, which we'll
    // deal with as an invalid.
    if cs < zcltok_first_final {
        f.emitToken(TokenInvalid, p, len(data))
    }

    // We always emit a synthetic EOF token at the end, since it gives the
    // parser position information for an "unexpected EOF" diagnostic.
    f.emitToken(TokenEOF, len(data), len(data))

    return f.Tokens
}
