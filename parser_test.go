package introduce_test

import (
	"reflect"
	"testing"

	"github.com/buildscale/introduce"
)

func TestParser(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		String   string
		Expected []introduce.ExpressionItem
	}{
		{
			String: `Buildscale... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}`,
			Expected: []introduce.ExpressionItem{
				{Text: "Buildscale... "},
				{Expansion: introduce.VariableExpansion{
					Identifier: "HELLO_WORLD",
				}},
				{Text: " "},
				{Expansion: introduce.EmptyValueExpansion{
					Identifier: "ANOTHER_VAR",
					Content: introduce.Expression([]introduce.ExpressionItem{{
						Text: "🏖",
					}}),
				}},
			},
		},
		{
			String: `${TEST1:- ${TEST2:-$TEST3}}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.EmptyValueExpansion{
					Identifier: "TEST1",
					Content: introduce.Expression([]introduce.ExpressionItem{
						{Text: " "},
						{Expansion: introduce.EmptyValueExpansion{
							Identifier: "TEST2",
							Content: introduce.Expression([]introduce.ExpressionItem{
								{Expansion: introduce.VariableExpansion{
									Identifier: "TEST3",
								}},
							}),
						}},
					}),
				}},
			},
		},
		{
			String: `${HELLO_WORLD-blah}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: introduce.Expression([]introduce.ExpressionItem{{
						Text: "blah",
					}}),
				}},
			},
		},
		{
			String: `\\${HELLO_WORLD-blah}`,
			Expected: []introduce.ExpressionItem{
				{Text: `\\`},
				{Expansion: introduce.UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: introduce.Expression([]introduce.ExpressionItem{{
						Text: "blah",
					}}),
				}},
			},
		},
		{
			String: `\${HELLO_WORLD-blah}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.EscapedExpansion{Identifier: "{HELLO_WORLD-blah}"}},
			},
		},
		{
			String: `Test \\\${HELLO_WORLD-blah}`,
			Expected: []introduce.ExpressionItem{
				{Text: `Test `},
				{Text: `\\`},
				{Expansion: introduce.EscapedExpansion{Identifier: "{HELLO_WORLD-blah}"}},
			},
		},
		{
			String: `${HELLO_WORLD:1}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
				}},
			},
		},
		{
			String: `${HELLO_WORLD: -1}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     -1,
				}},
			},
		},
		{
			String: `${HELLO_WORLD:-1}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.EmptyValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: introduce.Expression([]introduce.ExpressionItem{{
						Text: "1",
					}}),
				}},
			},
		},
		{
			String: `${HELLO_WORLD:1:7}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     7,
					HasLength:  true,
				}},
			},
		},
		{
			String: `${HELLO_WORLD:1:-7}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     -7,
					HasLength:  true,
				}},
			},
		},
		{
			String: `${HELLO_WORLD?Required}`,
			Expected: []introduce.ExpressionItem{
				{Expansion: introduce.RequiredExpansion{
					Identifier: "HELLO_WORLD",
					Message: introduce.Expression([]introduce.ExpressionItem{
						{Text: "Required"},
					}),
				}},
			},
		},
		{
			String: `$`,
			Expected: []introduce.ExpressionItem{
				{Text: `$`},
			},
		},
		{
			String: `\`,
			Expected: []introduce.ExpressionItem{
				{Text: `\`},
			},
		},
		{
			String: `$(echo hello world)`,
			Expected: []introduce.ExpressionItem{
				{Text: `$(`},
				{Text: `echo hello world)`},
			},
		},
		{
			String:   "$$MOUNTAIN",
			Expected: []introduce.ExpressionItem{{Expansion: introduce.EscapedExpansion{Identifier: "MOUNTAIN"}}},
		},
		{
			String: "this is a regex! /^start.*end$$/", // the dollar sign at the end of the regex has to be escaped to be treated as a literal dollar sign by this library
			Expected: []introduce.ExpressionItem{
				{Text: "this is a regex! /^start.*end"},
				{Expansion: introduce.EscapedExpansion{Identifier: ""}},
				{Text: "/"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.String, func(t *testing.T) {
			t.Parallel()

			actual, err := introduce.NewParser(tc.String).Parse()
			if err != nil {
				t.Fatal(err)
			}

			expected := introduce.Expression(tc.Expected)
			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Expected vs Actual: \n%s\n\n%s", expected, actual)
			}
		})
	}
}
